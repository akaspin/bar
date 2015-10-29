package command

import (
	"fmt"
	"github.com/akaspin/bar/client/git"
	"github.com/akaspin/bar/client/model"
	"github.com/spf13/cobra"
)

type GitDivertStatusCmd struct {
	*Environment
	*CommonOptions
}

func (c *GitDivertStatusCmd) Init(cc *cobra.Command) {
	cc.Use = "status"
	cc.Short = "covert op status"
}

func (c *GitDivertStatusCmd) Run(args ...string) (err error) {
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
		fmt.Fprintln(c.Stdout, "divert not in progress")
	}
	spec, err := divert.ReadSpec()
	if err != nil {
		return
	}
	fmt.Fprintln(c.Stdout, "DIVERT IN PROGRESS")
	fmt.Fprintln(c.Stdout, spec)

	return
}
