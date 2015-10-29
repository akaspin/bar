package command

import (
	"encoding/json"
	"fmt"
	"github.com/akaspin/bar/client/lists"
	"github.com/akaspin/bar/client/model"
	"github.com/akaspin/bar/client/transport"
	"github.com/akaspin/bar/proto"
	"github.com/spf13/cobra"
	"github.com/tamtam-im/logx"
	"path/filepath"
)

type SpecImportCmd struct {
	*Environment
	*CommonOptions

	UseGit bool
	Raw    bool
	Squash bool
}

func (c *SpecImportCmd) Init(cc *cobra.Command) {
	cc.Use = "import [spec-id]"
	cc.Short = "import spec from bard server"

	cc.Flags().BoolVarP(&c.UseGit, "git", "", false, "use git infrastructure")
	cc.Flags().BoolVarP(&c.Raw, "raw", "", false,
		"read spec from STDIN instead request from bar server")
	cc.Flags().BoolVarP(&c.Squash, "squash", "", false,
		"write manifests after import")
}

func (c *SpecImportCmd) Run(args ...string) (err error) {
	var spec proto.Spec

	mod, err := model.New(c.WD, c.UseGit, c.ChunkSize, c.PoolSize)
	if err != nil {
		return
	}
	trans := transport.NewTransport(mod, "", c.Endpoint, c.PoolSize)

	if c.Raw {
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

	idm := lists.IDMap{}
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

	if len(names) == 0 {
		logx.Fatalf("no manifests on server %s", names)
	}

	logx.Debugf("importing %s", names)

	if c.UseGit {
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

	if c.Squash {
		if err = mod.SquashBlobs(toSquash); err != nil {
			return
		}
	}
	for k, _ := range feed {
		fmt.Fprintf(c.Stdout, "%s ", filepath.FromSlash(k))
	}
	return
}
