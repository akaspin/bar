package command

import "github.com/spf13/cobra"

type GitDivertRootCmd struct {
}

func (c *GitDivertRootCmd) Init(cc *cobra.Command) {
	cc.Use = "divert"
	cc.Short = "git diversions"
}

func (c *GitDivertRootCmd) Run(args ...string) (err error) {
	return
}
