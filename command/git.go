package command

import "github.com/spf13/cobra"

// git root cmd
type GitCmd struct{}

func (c *GitCmd) Init(cc *cobra.Command) {
	cc.Use = "git"
	cc.Short = "git-specific operations"
}
func (c *GitCmd) Run(args ...string) (err error) { return }
