package cmd
import (
	"flag"
	"io"
	"github.com/akaspin/bar/barc/git"
	"os"
	"github.com/akaspin/bar/barc/lists"
	"github.com/akaspin/bar/shadow"
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
	useGit bool
	endpoint string
	noHeader bool
	fullID bool
	noRemote bool

	out io.Writer
	fs *flag.FlagSet
	wd string

	git *git.Git
	transport *transport.TransportPool
}

func (c *LsCmd) Bind(wd string, fs *flag.FlagSet, in io.Reader, out io.Writer) (err error) {
	fs.BoolVar(&c.useGit, "git", false, "use git infrastructure")
	fs.StringVar(&c.endpoint, "endpoint", "http://localhost:3000/v1",
		"bard endpoint")
	fs.BoolVar(&c.noHeader, "no-header", false, "do not show header")
	fs.BoolVar(&c.fullID, "full-id", false, "do not trim BLOB IDs")
	fs.BoolVar(&c.noRemote, "no-remote", false,
		"do not request bard for BLOBs states")

	c.fs = fs
	c.out = out
	c.wd = wd
	return
}

func (c *LsCmd) Do() (err error) {

	if c.useGit {
		if c.git, err = git.NewGit(""); err != nil {
			return
		}
	}

	// Collect filenames
	feed := lists.NewFileList(c.fs.Args()...).ListDir(c.wd)
	if c.git != nil {
		dirty, err := c.git.DiffFilesWithFilter(feed...)
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
	w.Init(c.out, 0, 8, 0, '\t', 0)

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
		Shadow *shadow.Shadow
		IsShadow bool},
	err error,
) {
	wg := &sync.WaitGroup{}
	res = map[string]struct{
		Shadow *shadow.Shadow
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
				Shadow *shadow.Shadow
				IsShadow bool
			}{res1, isSh}
		}(name)
	}


	wg.Wait()
	return
}

// Collect shadow for name
func (c *LsCmd) collectOneShadow(name string) (res *shadow.Shadow, isShadow bool, err error) {
	info, err := os.Stat(name)
	if err != nil {
		return
	}

	f, err := os.Open(name)
	if err != nil {
		return
	}
	defer f.Close()

	r, isShadow, err := shadow.Peek(f)
	if c.git == nil {
		// No git. Parse as usual
	} else {
		// If git is present - get manifest from index
		var oid string
		oid, err = c.git.OID(name)
		if err != nil {
			return
		}
		r, err = c.git.Cat(oid)
	}
	res, err = shadow.New(r, info.Size())
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
