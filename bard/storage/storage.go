package storage
import (
	"io"
)

type StorageFactory interface {
	GetStorage() (StorageDriver, error)
}

// All operations in storage driver are independent to each other
type StorageDriver interface {
	io.Closer

	IsExists(id string) (ok bool, err error)

	// Get BLOB shadow in full form
//	Describe(id string, out io.Writer) (err error)

	// Store BLOB from reader
	StoreBLOB(id string, size int64, in io.Reader) (err error)

	// Destroy blob
	DestroyBLOB(id string) (err error)
}


