package command
import (
	"github.com/spf13/cobra"
	"github.com/akaspin/bar/bar/model"
	"github.com/akaspin/bar/bar/transport"
	"github.com/tamtam-im/logx"
)

type PingCmd struct {
	*Environment
	*CommonOptions
}

func (c *PingCmd) Init(cc *cobra.Command)  {
	cc.Use = "ping"
	cc.Short = "ping bar server"
}

func (c *PingCmd) Run(args ...string) (err error) {
	var mod *model.Model
	if mod, err = model.New(c.WD, false, c.ChunkSize, c.PoolSize); err != nil {
		return
	}
	trans := transport.NewTransport(mod, "", c.Endpoint, c.PoolSize)
	res, err := trans.ServerInfo()
	logx.Info(res)
	return
}