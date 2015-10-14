package storage
import (
	"io"
	"github.com/akaspin/bar/manifest"
	"github.com/akaspin/bar/proto"
)

// All operations in storage driver are independent to each other
type Storage interface {
	io.Closer

	IsSpecExists(id manifest.ID) (ok bool, err error)

	IsBLOBExists(id manifest.ID) (ok bool, err error)

//	CheckBLOBS(ids []string) (map[string]bool, error)

	// Write spec
	WriteSpec(s proto.Spec) (err error)

	// Read spec
	ReadSpec(id manifest.ID) (res proto.Spec, err error)

	// Read manifest
	ReadManifest(id manifest.ID) (res *manifest.Manifest, err error)

	// Get manifests by it's ids
	GetManifests(ids []manifest.ID) (res []*manifest.Manifest, err error)

	// Declare new upload
	DeclareUpload(m manifest.Manifest) (err error)

	// Write chunk for declared blob from given reader
	WriteChunk(blobID, chunkID manifest.ID, size int64, r io.Reader) (err error)

	// Finish upload
	FinishUpload(id manifest.ID) (err error)

	// Read Chunk from blob by size and offset
	ReadChunkFromBlob(blobID manifest.ID, size, offset int64, w io.Writer) (err error)

//	GetManifests(ids [][]byte, )

}


