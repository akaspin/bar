package proto
import (
	"github.com/akaspin/bar/proto/manifest"
)

// Server info
type Info struct {

	// Alternate endpoints
	Endpoints []string

	// Preferred chunk size
	ChunkSize int64

	// Preferred number of connections from client
	MaxConn int
}

type ChunkInfo struct {
	BlobID string
	manifest.Chunk
}

type ChunkData struct {
	ChunkInfo
	Data []byte
}

