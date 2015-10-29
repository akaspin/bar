package command

import (
	"github.com/akaspin/bar/proto"
	"github.com/akaspin/bar/server"
	"github.com/akaspin/bar/server/front"
	"github.com/akaspin/bar/server/storage"
	"github.com/akaspin/bar/server/thrift"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"strings"
)

type ServerRunCmd struct {
	*Environment
	*CommonOptions
	*ServerOptions
}

func (c *ServerRunCmd) Init(cc *cobra.Command) {
	cc.Use = "run"
	cc.Short = "run bar server instance (bard)"

	cc.Flags().StringVarP(&c.BindHTTP, "bind-http", "", ":3000", "http bind")
	cc.Flags().StringVarP(&c.BindRPC, "bind-rpc", "", ":3001", "rpc bind")

	cc.Flags().StringVarP(&c.HTTPEndpoint, "endpoint-http", "",
		"http://localhost:3000/v1", "http endpoint")
	cc.Flags().StringVarP(&c.BinDir, "bin-dir", "", "dist/bindir",
		"binaries directory")
	cc.Flags().StringVarP(&c.Storage, "storage", "",
		"block:root=data,split=2,max-files=128,pool=64",
		"storage configuration")
}

func (c *ServerRunCmd) Run(args ...string) (err error) {
	stor, err := storage.GuessStorage(c.Storage)
	if err != nil {
		return
	}

	info := &proto.ServerInfo{
		HTTPEndpoint: c.ServerOptions.HTTPEndpoint,
		RPCEndpoints: strings.Split(c.Endpoint, ","),
		ChunkSize:    c.ChunkSize,
		PoolSize:     c.PoolSize,
		BufferSize:   c.BufferSize,
	}
	ctx := context.Background()
	// Thrift
	tServer := thrift.NewServer(ctx, &thrift.Options{info, c.BindRPC}, stor)
	httpServer := front.NewServer(ctx, &front.Options{info, c.BindHTTP, c.BinDir}, stor)

	srv := server.NewCompositeServer(ctx, tServer, httpServer)

	err = srv.Start()

	return
}
