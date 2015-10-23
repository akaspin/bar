package command
import (
	"github.com/spf13/cobra"
	"github.com/akaspin/bar/bar/model"
	"github.com/akaspin/bar/bar/git"
	"fmt"
)

type GitDivertFinishCmd struct {
	*Environment
	*CommonOptions

	Message string
}

func (c *GitDivertFinishCmd) Init(cc *cobra.Command) {
	cc.Use = "finish"
	cc.Short = "finish covert op"

	cc.Flags().StringVarP(&c.Message, "message", "m", "",
		"git commit message")
}

func (c *GitDivertFinishCmd) Run(args ...string) (err error) {
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
	if err = divert.Commit(spec, c.Message); err != nil {
		return
	}

	err = divert.Cleanup(spec)


	return
}