package server

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