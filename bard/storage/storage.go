package storage
import (
	"io"
	"github.com/akaspin/bar/proto/manifest"
	"github.com/akaspin/bar/proto"
)

type StorageFactory interface {
	GetStorage() (StorageDriver, error)
}

// All operations in storage driver are independent to each other
type StorageDriver interface {
	io.Closer

	IsExists(id string) (ok bool, err error)

	// Write spec
	WriteSpec(s proto.Spec) (err error)

	// Read spec
	ReadSpec(id string) (res proto.Spec, err error)

	// Read manifest
	ReadManifest(id string) (res manifest.Manifest, err error)

	// Declare new upload
	DeclareUpload(m manifest.Manifest) (err error)

	// Write chunk for declared blob
	WriteChunk(blobID, chunkID string, size int64, r io.Reader) (err error)

	// Finish upload
	FinishUpload(id string) (err error)

	// Destroy blob
	DestroyBLOB(id string) (err error)

	// Get BLOB stream
	ReadBLOB(id string) (r io.ReadCloser, err error)
}


