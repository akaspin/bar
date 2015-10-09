package manifest
import (
	"github.com/vireshas/minimal_vitess_pool/pools"
	"time"
	"io"
)

type Hasher struct {
	chunkSize int64
	pool *HasherPool
}

func (h *Hasher) Make(in io.Reader) (res *Manifest, err error) {
	res, err = NewFromAny(in, h.chunkSize)
	return
}

func (h *Hasher) MakeFromManifest(in io.Reader) (res *Manifest, err error) {
	res, err = NewFromManifest(in)
	return
}

func (h *Hasher) MakeFromBLOB(in io.Reader) (res *Manifest, err error) {
	res, err = NewFromBLOB(in, h.chunkSize)
	return
}

func (h *Hasher) Peek(in io.Reader) (r io.Reader, isManifest bool, err error) {
	return Peek(in)
}

func (h *Hasher) Release() {
	h.pool.pool.Put(h)
}

func (h *Hasher) Close() {}


// Concurrent hashing
type HasherPool struct {
	pool *pools.ResourcePool
}

func NewHasherPool(chunkSize int64, n int, timeout time.Duration) (res *HasherPool) {
	res = &HasherPool{}
	newFn := func() (pools.Resource, error) {
		return &Hasher{chunkSize, res}, nil
	}
	return &HasherPool{pools.NewResourcePool(newFn, n, n, timeout)}
}

func (p *HasherPool) Take() (res *Hasher, err error) {
	h1, err := p.pool.Get(time.Minute * 30)
	if err != nil {
		return
	}
	res = h1.(*Hasher)
	return
}

// Make from reader
func (p *HasherPool) Make(in io.Reader) (res *Manifest, err error)  {
	h1, err := p.pool.Get(time.Minute * 30)
	if err != nil {
		return
	}
	defer p.pool.Put(h1)
	h := h1.(*Hasher)

	res, err = h.Make(in)
	return
}

