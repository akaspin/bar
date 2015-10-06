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

	// Write manifest
	WriteManifest(m manifest.Manifest) (err error)

	// Read manifest
	ReadManifest(id string) (res manifest.Manifest, err error)

	// Store BLOB from given reader
	WriteBLOB(id string, size int64, in io.Reader) (err error)

	// Destroy blob
	DestroyBLOB(id string) (err error)

	// Get BLOB stream
	ReadBLOB(id string) (r io.ReadCloser, err error)
}


