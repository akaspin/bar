package command
import (
	"github.com/spf13/cobra"
	"fmt"
)

var Version string

type VersionCmd struct {
	*Environment
}

func (c *VersionCmd) Init(cc *cobra.Command) {
	cc.Use = "version"
	cc.Short = "print version and exit"
}

func (c *VersionCmd) Run(args ...string) (err error)  {
	fmt.Fprint(c.Stdout, Version)
	return
}