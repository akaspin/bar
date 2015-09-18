package storage
import (
	"github.com/vireshas/minimal_vitess_pool/pools"
	"time"
)


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
