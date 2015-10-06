package manifest
import (
	"github.com/vireshas/minimal_vitess_pool/pools"
	"time"
	"io"
)

type hasher struct {
	chunkSize int64
}

func (h *hasher) Make(in io.Reader) (res *Manifest, err error) {
	res, err = NewFromAny(in, h.chunkSize)
	return
}

func (h *hasher) MakeFromManifest(in io.Reader) (res *Manifest, err error) {
	res, err = NewFromManifest(in)
	return
}

func (h *hasher) Close() {}


// Concurrent hashing
type Hasher struct {
	pool *pools.ResourcePool
}

func NewHasherPool(chunkSize int64, n int, timeout time.Duration) *Hasher {
	newFn := func() (pools.Resource, error) {
		return &hasher{chunkSize}, nil
	}
	return &Hasher{pools.NewResourcePool(newFn, n, n, timeout)}
}

// Make from reader
func (p *Hasher) Make(in io.Reader) (res *Manifest, err error)  {
	h1, err := p.pool.TryGet()
	if err != nil {
		return
	}
	defer p.pool.Put(h1)
	h := h1.(*hasher)

	res, err = h.Make(in)
	return
}

