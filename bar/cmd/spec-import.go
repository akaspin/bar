package cmd
import (
	"flag"
	"github.com/akaspin/bar/bar/model"
	"github.com/akaspin/bar/bar/transport"
	"encoding/json"
	"github.com/akaspin/bar/proto"
	"github.com/tamtam-im/logx"
	"github.com/akaspin/bar/bar/lists"
	"fmt"
	"path/filepath"
)


/*
Import spec from bard and populate manifests

	$ cat spec.json | bar spec-import -raw
	$ bar spec-import 1bcaa5...578bd24
	$ bar spec-import http://localhost:3000/v1/spec/1bcaa5...578bd24
	$ bar spec-import http://localhost:3000/v1/spec/1bcaa5...578bd24.json

*/
type SpecImportCmd struct  {
	*Base

	useGit bool
	raw bool
	squash bool
}

func NewSpecImportCmd(s *Base) SubCommand  {
	c := &SpecImportCmd{Base: s}
	return c
}

func (c *SpecImportCmd) Init(fs *flag.FlagSet) {
	fs.BoolVar(&c.useGit, "git", false, "use git infrastructure")
	fs.BoolVar(&c.raw, "raw", false, "read spec from STDIN")
	fs.BoolVar(&c.squash, "squash", false, "squash modified blobs")
}

func (c *SpecImportCmd) Do(args []string) (err error) {
	// get spec first
	var spec proto.Spec

	mod, err := model.New(c.WD, c.useGit, c.ChunkSize, c.PoolSize)
	if err != nil {
		return
	}
	trans := transport.NewTransport(mod, "", c.Endpoints, c.PoolSize)

	if c.raw {
		if err = json.NewDecoder(c.Stdin).Decode(&spec); err != nil {
			return
		}
	} else {
		// tree spec types
		id := proto.ID(args[0])

		if spec, err = trans.GetSpec(id); err != nil {
			logx.Debug(spec, err)
			return
		}
	}

	var idm lists.IDMap
	for n, id := range spec.BLOBs {
		idm[id] = append(idm[id], n)
	}

	// request manifests and
	mans, err := trans.GetManifests(idm.IDs())
	if err != nil {
		return
	}
	feed := idm.ToBlobMap(mans)
	names := feed.Names()

	logx.Debugf("importing %s", names)

	if c.useGit {
		// If git is used - check names for attrs
		byAttr, err := mod.Git.FilterByAttr("bar", names...)
		if err != nil {
			return err
		}

		diff := []string{}
		attrs := map[string]struct{}{}
		for _, x := range byAttr {
			attrs[x] = struct{}{}
		}

		for _, x := range names {
			if _, ok := attrs[x]; !ok {
				diff = append(diff, x)
			}
		}
		if len(diff) > 0 {
			return fmt.Errorf("some spec blobs is not under bar control %s", diff)
		}
	}

	// get stored links, ignore errors
	stored, _ := mod.FeedManifests(true, true, false, names...)

	logx.Debugf("already stored %s", stored.Names())

	// squash present
	toSquash := lists.BlobMap{}
	for n, m := range feed {
		m1, ok := stored[filepath.FromSlash(n)]
		if !ok || m.ID != m1.ID {
			toSquash[n] = feed[n]
		}
	}

	if c.squash {
		if err = mod.SquashBlobs(toSquash); err != nil {
			return
		}
	}
	for k, _ := range feed {
		fmt.Fprintf(c.Stdout, "%s ", filepath.FromSlash(k))
	}
	return
}

func (c *SpecImportCmd) Help() {
	fmt.Fprintln(c.Stderr, "bar spec-import [OPTIONS] [SPEC-ID]\n")
}

func (c *SpecImportCmd) Description() string {
	return "import spec from bar server"
}

