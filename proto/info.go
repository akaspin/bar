package proto

import (
	"github.com/akaspin/bar/proto/wire"
	"strings"
)

// Server info
type ServerInfo struct {

	// HTTP Endpoint
	HTTPEndpoint string

	// RPC Endpoint
	RPCEndpoints []string

	// Preferred chunk size
	ChunkSize int64

	// Preferred number of connections from client
	PoolSize int

	// Thrift rpc buffer size
	BufferSize int
}

func (i *ServerInfo) JoinRPCEndpoints() string {
	return strings.Join(i.RPCEndpoints, ",")
}

func (i ServerInfo) MarshalThrift() (res wire.ServerInfo, err error) {
	res = wire.ServerInfo{
		HttpEndpoint: i.HTTPEndpoint,
		RpcEndpoints: i.RPCEndpoints,
		ChunkSize:    i.ChunkSize,
		MaxConn:      int32(i.PoolSize),
		BufferSize:   int32(i.BufferSize),
	}
	return
}

func (i *ServerInfo) UnmarshalThrift(data wire.ServerInfo) (err error) {
	i.HTTPEndpoint = data.HttpEndpoint
	i.RPCEndpoints = data.RpcEndpoints
	i.ChunkSize = data.ChunkSize
	i.PoolSize = int(data.MaxConn)
	i.BufferSize = int(data.BufferSize)
	return
}
