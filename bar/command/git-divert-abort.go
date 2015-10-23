package command
import (
	"github.com/spf13/cobra"
	"github.com/akaspin/bar/bar/model"
	"github.com/akaspin/bar/bar/git"
	"fmt"
)

type GitDivertAbortCmd struct  {
	*Environment
	*CommonOptions
}

func (c *GitDivertAbortCmd) Init(cc *cobra.Command) {
	cc.Use = "abort"
	cc.Short = "abort covert op"
}

func (c *GitDivertAbortCmd) Run(args ...string) (err error) {
	mod, err := model.New(c.WD, true, c.ChunkSize, c.PoolSize)
	if err != nil {
		return
	}

	divert := git.NewDivert(mod.Git)
	isInProgress, err := divert.IsInProgress()
	if err != nil {
		return
	}
	if !isInProgress {
		err = fmt.Errorf("diversion is not in progress")
	}

	spec, err := divert.ReadSpec()
	if err != nil {
		return
	}
	err = divert.Cleanup(spec)
	return
}