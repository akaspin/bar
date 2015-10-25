package command
import (
	"github.com/spf13/cobra"
	"github.com/akaspin/bar/client/model"
	"github.com/akaspin/bar/client/lists"
	"fmt"
	"github.com/tamtam-im/logx"
	"github.com/akaspin/bar/client/transport"
)

type UpCmd struct  {
	*Environment
	*CommonOptions

	UseGit bool
	Squash bool
}

func (c *UpCmd) Init(cc *cobra.Command) {
	cc.Use = "up"
	cc.Short = "upload BLOBs to bar server"
	cc.Flags().BoolVarP(&c.UseGit, "git", "", false, "use git infrastructure")
	cc.Flags().BoolVarP(&c.Squash, "squash", "", false,
		"squash uploaded BLOBs to manifests after upload")
}

func (c *UpCmd) Run(args ...string) (err error) {
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

	blobs, err := mod.FeedManifests(true, false, true, feed...)
	if err != nil {
		return
	}

	logx.Debugf("collected blobs %s", blobs.IDMap())

	trans := transport.NewTransport(mod, "", c.Endpoint, c.PoolSize)

	err = trans.Upload(blobs)
	if err != nil {
		return
	}

	if c.Squash {
		if err = mod.SquashBlobs(blobs); err != nil {
			return
		}
		if c.UseGit {
			err = mod.Git.UpdateIndex(blobs.Names()...)
		}
	}

	return
}

