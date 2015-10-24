package command
import (
	"github.com/akaspin/bar/server"
	"github.com/spf13/cobra"
	"github.com/akaspin/bar/bard/storage"
)

type ServerRunCmd struct  {
	*Environment
	*CommonOptions
	*server.ServerOptions
}

func (c *ServerRunCmd) Init(cc *cobra.Command) {
	cc.Use = "run"
	cc.Short = "run bar server instance (bard)"

	cc.Flags().StringVarP(&c.BindHTTP, "bind-http", "", ":3000", "http bind")
	cc.Flags().StringVarP(&c.BindRPC, "bind-rpc", "", ":3001", "rpc bind")

	cc.Flags().StringVarP(&c.HTTPEndpoint, "endpoint-http",
		"http://localhost:3000/v1", "http endpoint")
	cc.Flags().StringVarP(&c.BinDir, "bin-dir", "", "binaries directory")
	cc.Flags().StringVarP(&c.Storage, "storage", "",
		"block:root=data,split=2,max-files=128,pool=64",
		"storage configuration")
}

func (c *ServerRunCmd) Run(args ...string) (err error) {

	store, err := storage.GuessStorage(c.Storage)
	if err != nil {
		return
	}



	return
}