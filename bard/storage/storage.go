package storage
import (
	"io"
	"github.com/vireshas/minimal_vitess_pool/pools"
	"time"
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


type poolWrapper struct {
	s StorageDriver
}

func (w *poolWrapper) Close() {
	w.s.Close()
}

type StoragePool struct  {
	p *pools.ResourcePool
}

func NewStoragePool(factory StorageFactory, max int, timeout time.Duration) *StoragePool {
	newFn := func() (res pools.Resource, err error) {
		s, err := factory.GetStorage()
		if err != nil {
			return
		}
		res = &poolWrapper{s}
		return
	}
	return &StoragePool{pools.NewResourcePool(newFn, max, max, timeout)}
}

func (p *StoragePool) Take() (res StorageDriver, err error) {
	r, err := p.p.TryGet()
	if err != nil {
		return
	}
	res = r.(*poolWrapper).s
	return
}

func (p *StoragePool) Release(s StorageDriver) {
	p.p.Put(&poolWrapper{s})
}
