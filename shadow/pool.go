package shadow
import (
	"github.com/vireshas/minimal_vitess_pool/pools"
	"time"
	"io"
)

type hasher struct {
	chunkSize int64
}

func (h *hasher) Shadow(in io.Reader, full bool) (res *Shadow, err error) {
	res = &Shadow{}
	err = res.FromAny(in, full, h.chunkSize)
	return
}

type hasherResource struct {
	h *hasher
}

func (w *hasherResource) Close() {}

// Concurrent hashing
type HasherPool struct {
	pool *pools.ResourcePool
}

func NewHasherPool(n int, timeout time.Duration, chunkSize int64) *HasherPool {
	newFn := func() (pools.Resource, error) {
		return &hasherResource{&hasher{chunkSize}}, nil
	}
	return &HasherPool{pools.NewResourcePool(newFn, n, n, timeout)}
}

// Make one shadow from reader
func (p *HasherPool) MakeOne(in io.Reader, full bool) (res *Shadow, err error)  {
	r, err := p.pool.TryGet()
	if err != nil {
		return
	}
	defer p.pool.Put(r)
	h := r.(*hasherResource).h
	res, err = h.Shadow(in, full)
	return
}
