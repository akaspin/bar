package shadow
import (
	"github.com/vireshas/minimal_vitess_pool/pools"
	"time"
	"io"
)

type hasher struct {
}

func (h *hasher) Shadow(in io.Reader, size int64) (res *Shadow, err error) {
	res = &Shadow{}
	res, err = New(in, size)
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

func NewHasherPool(n int, timeout time.Duration) *HasherPool {
	newFn := func() (pools.Resource, error) {
		return &hasherResource{&hasher{}}, nil
	}
	return &HasherPool{pools.NewResourcePool(newFn, n, n, timeout)}
}

// Make one shadow from reader
func (p *HasherPool) MakeOne(in io.Reader, size int64) (res *Shadow, err error)  {
	r, err := p.pool.TryGet()
	if err != nil {
		return
	}
	defer p.pool.Put(r)
	h := r.(*hasherResource).h
	res, err = h.Shadow(in, size)
	return
}
