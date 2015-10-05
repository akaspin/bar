package cmd
import (
	"github.com/akaspin/bar/barc/git"
	"os"
	"github.com/akaspin/bar/barc/lists"
	"github.com/akaspin/bar/proto/manifest"
	"sync"
	"fmt"
	"text/tabwriter"
	"github.com/akaspin/bar/barc/transport"
	"net/url"
	"time"
)

/*
List bar blobs

	$ barc ls my
	NAME                SHADOW   SYNC   ID              SIZE
	my/blob             -        +      3d3463063cb1    1930746935
	my/other/blob       +        +      309a34901901    3
	my/bad/blob         -        -      3d3463063cb1    1930746935


*/
type LsCmd struct {
	*BaseSubCommand

	useGit bool
	endpoint string
	noHeader bool
	fullID bool
	noRemote bool

	git *git.Git
	transport *transport.TransportPool
}

func NewLsCmd(s *BaseSubCommand) SubCommand {
	c := &LsCmd{BaseSubCommand: s}
	c.FS.BoolVar(&c.useGit, "git", false, "use git infrastructure")
	c.FS.StringVar(&c.endpoint, "endpoint", "http://localhost:3000/v1",
		"bard endpoint")
	c.FS.BoolVar(&c.noHeader, "no-header", false, "do not show header")
	c.FS.BoolVar(&c.fullID, "full-id", false, "do not trim BLOB IDs")
	c.FS.BoolVar(&c.noRemote, "no-remote", false,
		"do not request bard for BLOBs states")
	return c
}

func (c *LsCmd) Do() (err error) {

	if c.useGit {
		if c.git, err = git.NewGit(""); err != nil {
			return
		}
	}

	// Collect filenames
	feed := lists.NewFileList(c.FS.Args()...).ListDir(c.WD)
	if c.git != nil {
		dirty, err := c.git.DiffFilesWithAttr(feed...)
		if err != nil {
			return err
		}
		if len(dirty) > 0 {
			return fmt.Errorf("dirty files in tree %s", dirty)
		}
	}
	shadows, err := c.collectShadows(feed)

	// request bard
	ids := []string{}
	for _, s := range shadows {
		ids = append(ids, s.Shadow.ID)
	}

	remotes, err := c.remote(ids)
	if err != nil {
		return
	}

	// print this stuff
	w := new(tabwriter.Writer)
	w.Init(c.Stdout, 0, 8, 0, '\t', 0)

	if !c.noHeader {
		fmt.Fprint(w, "NAME\tBLOB\t")
		if !c.noRemote {
			fmt.Fprint(w, "REMOTE\t")
		}
		fmt.Fprintln(w, "ID\tSIZE")
	}

	for _, name := range feed {
		sh := shadows[name]
		blob := "yes"
		if sh.IsShadow {
			blob = "no"
		}
		fmt.Fprintf(w, "%s\t%s\t", name, blob)

		if !c.noRemote {
			remote := "no"
			if _, ok := remotes[sh.Shadow.ID]; ok {
				remote = "yes"
			}
			fmt.Fprintf(w, "%s\t", remote)
		}
		id := sh.Shadow.ID
		if !c.fullID {
			id = sh.Shadow.ID[:12]
		}
		fmt.Fprintf(w, "%s\t%d", id, sh.Shadow.Size)
		fmt.Fprintln(w)
	}
	w.Flush()

	return
}

// Collect all shadows for feed
func (c *LsCmd) collectShadows(feed []string) (
	res map[string]struct{
		Shadow *manifest.Manifest
		IsShadow bool},
	err error,
) {
	wg := &sync.WaitGroup{}
	res = map[string]struct{
		Shadow *manifest.Manifest
		IsShadow bool
	}{}
	errs := []error{}
	for _, name := range feed {
		wg.Add(1)
		go func(n string) {
			defer wg.Done()
			res1, isSh, err1 := c.collectOneShadow(n)
			if err1 != nil {
				errs = append(errs, err1)
				return
			}
			res[n] = struct{
				Shadow *manifest.Manifest
				IsShadow bool
			}{res1, isSh}
		}(name)
	}


	wg.Wait()
	return
}

// Collect shadow for name
func (c *LsCmd) collectOneShadow(name string) (res *manifest.Manifest, isShadow bool, err error) {

	f, err := os.Open(name)
	if err != nil {
		return
	}
	defer f.Close()

	r, isShadow, err := manifest.Peek(f)
	if c.git == nil {
		// No git. Parse as usual
	} else {
		// If git is present - get manifest from index
		var oid string
		oid, err = c.git.GetOID(name)
		if err != nil {
			return
		}
		r, err = c.git.Cat(oid)
	}
	res, err = manifest.NewFromAny(r, manifest.CHUNK_SIZE)
	return
}

func (c *LsCmd) remote(in []string) (res map[string]struct{}, err error) {
	res = map[string]struct{}{}
	if c.noRemote {
		return
	}
	u, err := url.Parse(c.endpoint)
	if err != nil {
		return
	}

	c.transport = transport.NewTransportPool(u, 10, time.Minute)
	tr, err := c.transport.Take()
	if err != nil {
		return
	}
	defer c.transport.Release(tr)

	resp, err := tr.Check(in)
	if err != nil {
		return
	}
	for _, id := range resp {
		res[id] = struct{}{}
	}
	return
}
