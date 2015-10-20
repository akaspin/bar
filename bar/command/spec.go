package command
import "github.com/spf13/cobra"

type SpecRootCmd struct {}

func (c *SpecRootCmd) Init(cc *cobra.Command)  {
	cc.Use = "spec"
	cc.Short = "spec operations"
}

func (c *SpecRootCmd) Run(args ...string) (err error)  {
	return
}