package command
import "github.com/spf13/cobra"

type GitDivertAbortCmd struct  {
	*Environment
	*CommonOptions
}

func (c *GitDivertAbortCmd) Init(cc *cobra.Command) {
	cc.Use = "abort"
	cc.Short = "abort covert op"
}

func (c *GitDivertAbortCmd) Run(args ...string) (err error) {
	return
}