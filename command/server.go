package command
import "github.com/spf13/cobra"

type ServerCmd struct {

}

func (c *ServerCmd) Init(cc *cobra.Command)  {
	cc.Use = "server"
	cc.Short = "bar server"
}

func (c *ServerCmd) Run(args ...string) (err error) {
	return
}