package command
import (
	"github.com/spf13/cobra"
	"github.com/akaspin/bar/client/model"
	"fmt"
	"github.com/akaspin/bar/client/transport"
	"github.com/akaspin/bar/client/git"
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
	var filenames []string
	var mod *model.Model
	if mod, err = model.New(c.WD, true, c.ChunkSize, c.PoolSize); err != nil {
		return
	}

	// In divert we need restrict check by target filenames
	divert := git.NewDivert(mod.Git)
	isInDivert, err := divert.IsInProgress()
	if err != nil {
		return
	}
	if isInDivert {
		var spec git.DivertSpec
		if spec, err = divert.ReadSpec(); err != nil {
			return
		}
		filenames = spec.TargetFiles
	}

	isDirty, dirty, err := mod.Check(filenames...)
	if err != nil {
		return
	}
	if isDirty {
		err = fmt.Errorf("dirty files in working tree %s", dirty)
		return
	}

	feedR, err := mod.Git.Diff(filenames...)
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