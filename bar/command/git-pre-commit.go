package command
import (
	"github.com/spf13/cobra"
	"github.com/akaspin/bar/bar/model"
	"fmt"
"github.com/akaspin/bar/bar/transport"
)

type GitPreCommitCmd struct {
	*Environment
	*CommonOptions
}

func (c *GitPreCommitCmd) Init(cc *cobra.Command) {
	cc.Use = "pre-commit"
	cc.Short = "git pre-commit hook"
}

func (c *GitPreCommitCmd) Run(args ...string) (err error)  {
	var mod *model.Model
	if mod, err = model.New(c.WD, true, c.ChunkSize, c.PoolSize); err != nil {
		return
	}

	isDirty, dirty, err := mod.Check()
	if err != nil {
		return
	}
	if isDirty {
		err = fmt.Errorf("dirty files in working tree %s", dirty)
		return
	}

	feedR, err := mod.Git.Diff()
	if err != nil {
		return
	}

	blobs, err := mod.Git.ManifestsFromDiff(feedR)
	if err != nil {
		return
	}

	trans := transport.NewTransport(mod, "", c.Endpoint, c.PoolSize)
	err = trans.Upload(blobs)

	return
}