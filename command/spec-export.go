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
	"os"
	"path/filepath"
	"time"
)

type SpecExportCmd struct {
	*Environment
	*CommonOptions

	UseGit bool

	Upload bool
	DoCC   bool
}

func (c *SpecExportCmd) Init(cc *cobra.Command) {
	cc.Use = "export [# path]"
	cc.Short = "export spec"

	cc.Flags().BoolVarP(&c.UseGit, "git", "", false, "use git infrastructure")
	cc.Flags().BoolVarP(&c.Upload, "upload", "u", false, "upload spec to bar server")
	cc.Flags().BoolVarP(&c.DoCC, "cc", "", false, "create carbon copy")
}

func (c *SpecExportCmd) Run(args ...string) (err error) {
	var mod *model.Model
	if mod, err = model.New(c.WD, c.UseGit, c.ChunkSize, c.PoolSize); err != nil {
		return
	}

	feed := lists.NewFileList(args...).ListDir(c.WD)

	isDirty, dirty, err := mod.Check(feed...)
	if err != nil {
		return
	}
	if isDirty {
		err = fmt.Errorf("dirty files in working tree %s", dirty)
		return
	}

	if c.UseGit {
		// filter by attrs
		feed, err = mod.Git.FilterByAttr("bar", feed...)
	}

	blobs, err := mod.FeedManifests(true, true, true, feed...)
	if err != nil {
		return
	}

	// make specmap
	nameMap := map[string]proto.ID{}
	for name, m := range blobs {
		nameMap[name] = m.ID
	}

	spec, err := proto.NewSpec(time.Now().UnixNano(), nameMap, []string{})
	if err != nil {
		return
	}

	if c.DoCC {
		ccName := fmt.Sprintf("bar-spec-%d-%s.json",
			time.Now().UnixNano(), spec.ID)
		logx.Infof("storing carbon copy to %s", ccName)
		ccf, err := os.Create(filepath.Join(c.WD, ccName))
		if err != nil {
			return err
		}
		defer ccf.Close()

		if err = json.NewEncoder(ccf).Encode(&spec); err != nil {
			return err
		}
	}

	if !c.Upload {
		err = json.NewEncoder(c.Stdout).Encode(&spec)
		return
	}

	trans := transport.NewTransport(mod, "", c.Endpoint, c.PoolSize)
	if err = trans.UploadSpec(spec); err != nil {
		return
	}
	fmt.Fprint(c.Stdout, spec.ID)
	return
}
