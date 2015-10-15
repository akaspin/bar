package proto
import (
	"strings"
)

// Server info
type Info struct {

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

func (i *Info) JoinRPCEndpoints() string {
	return strings.Join(i.RPCEndpoints, ",")
}

type ChunkInfo struct {
	BlobID ID
	Chunk
}

type ChunkData struct {
	ChunkInfo
	Data []byte
}

