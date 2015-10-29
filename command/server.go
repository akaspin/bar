package command

import "github.com/spf13/cobra"

type ServerOptions struct {
	// Http bind `:3000`
	BindHTTP string

	// rpc bind `:3001`
	BindRPC string

	// HTTP endpoint
	HTTPEndpoint string

	// Binaries directory
	BinDir string

	// Storage init
	//
	//   <storage-type>:option=...,option=...
	//   block:root=data,split=2,max-files=32,pool=32
	//
	Storage string
}

type ServerCmd struct {
}

func (c *ServerCmd) Init(cc *cobra.Command) {
	cc.Use = "server"
	cc.Short = "bar server"
}

func (c *ServerCmd) Run(args ...string) (err error) {
	return
}
