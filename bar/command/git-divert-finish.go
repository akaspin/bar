package command
import "github.com/spf13/cobra"

type GitDivertFinishCmd struct {
	*Environment
	*CommonOptions
}

func (c *GitDivertFinishCmd) Init(cc *cobra.Command) {
	cc.Use = "finish"
	cc.Short = "finish covert op"
}

func (c *GitDivertFinishCmd) Run(args ...string) (err error) {
	return
}