package command
import (
	"github.com/spf13/cobra"
	"github.com/akaspin/bar/client/model"
	"github.com/akaspin/bar/client/lists"
	"fmt"
	"github.com/akaspin/bar/client/transport"
)

type DownCmd struct  {
	*Environment
	*CommonOptions

	UseGit bool
}

func (c *DownCmd) Init(cc *cobra.Command) {
	cc.Use = "down"
	cc.Short = "download BLOBs from bar server"

	cc.Flags().BoolVarP(&c.UseGit, "git", "", false, "use git infrastructure")
}

func (c *DownCmd) Run(args ...string) (err error) {
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

	blobs, err := mod.FeedManifests(false, true, true, feed...)
	if err != nil {
		return
	}

	trans := transport.NewTransport(mod, "", c.Endpoint, c.PoolSize)
	if err = trans.Download(blobs); err != nil {
		return
	}

	if c.UseGit {
		err = mod.Git.UpdateIndex(blobs.Names()...)
	}

	return
}