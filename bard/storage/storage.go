package storage
import "io"

type Storage interface {

	// Store BLOB from reader
	Store(id []byte, size int64, in io.Reader) (err error)

	// Destroy blob
	Destroy(id string) (err error)
}
