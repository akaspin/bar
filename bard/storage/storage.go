package storage
import (
	"io"
	"github.com/akaspin/bar/proto"
)

// All operations in storage driver are independent to each other
type Storage interface {
	io.Closer

	IsSpecExists(id proto.ID) (ok bool, err error)

	IsBLOBExists(id proto.ID) (ok bool, err error)

//	CheckBLOBS(ids []string) (map[string]bool, error)

	// Write spec
	WriteSpec(s proto.Spec) (err error)

	// Read spec
	ReadSpec(id proto.ID) (res proto.Spec, err error)

	// Read proto
	ReadManifest(id proto.ID) (res *proto.Manifest, err error)

	// Get manifests by it's ids
	GetManifests(ids []proto.ID) (res []proto.Manifest, err error)

	// Declare new upload
	DeclareUpload(m proto.Manifest) (err error)

	// Write chunk for declared blob from given reader
	WriteChunk(blobID, chunkID proto.ID, size int64, r io.Reader) (err error)

	// Finish upload
	FinishUpload(id proto.ID) (err error)

	// Read Chunk from blob by size and offset
	ReadChunkFromBlob(blobID proto.ID, size, offset int64, w io.Writer) (err error)

//	GetManifests(ids [][]byte, )

}


