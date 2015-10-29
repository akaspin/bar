package command

import (
	"github.com/akaspin/bar/client/model"
	"github.com/akaspin/bar/client/transport"
	"github.com/spf13/cobra"
	"github.com/tamtam-im/logx"
)

type PingCmd struct {
	*Environment
	*CommonOptions
}

func (c *PingCmd) Init(cc *cobra.Command) {
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
