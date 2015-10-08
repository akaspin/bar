package proto
import (
	"github.com/akaspin/bar/proto/manifest"
)

// Server info
type Info struct {

	// HTTP Endpoint
	HTTPEndpoint string

	// RPC Endpoint
	Endpoint string

	// Preferred chunk size
	ChunkSize int64

	// Preferred number of connections from client
	PoolSize int
}

type ChunkInfo struct {
	BlobID string
	manifest.Chunk
}

type ChunkData struct {
	ChunkInfo
	Data []byte
}

