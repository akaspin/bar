package storage
import (
	"io"
	"github.com/akaspin/bar/proto/manifest"
	"github.com/akaspin/bar/proto"
)

// All operations in storage driver are independent to each other
type Storage interface {
	io.Closer

	IsSpecExists(id string) (ok bool, err error)

	IsBLOBExists(id string) (ok bool, err error)

//	CheckBLOBS(ids []string) (map[string]bool, error)

	// Write spec
	WriteSpec(s proto.Spec) (err error)

	// Read spec
	ReadSpec(id string) (res proto.Spec, err error)

	// Read manifest
	ReadManifest(id string) (res manifest.Manifest, err error)

	// Declare new upload
	DeclareUpload(m manifest.Manifest) (err error)

	// Write chunk for declared blob from given reader
	WriteChunk(blobID, chunkID string, size int64, r io.Reader) (err error)

	// Finish upload
	FinishUpload(id string) (err error)

	// Read chunk from storage to given writer
	ReadChunk(chunk proto.ChunkInfo, w io.Writer) (err error)

	// Read Chunk from blob by size and offset
	ReadChunkFromBlob(blobID []byte, size, offset int64, w io.Writer) (err error)

	// Destroy blob
	DestroyBLOB(id string) (err error)
}


