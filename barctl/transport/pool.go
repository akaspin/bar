package transport
import (
	"github.com/vireshas/minimal_vitess_pool/pools"
	"net/url"
	"time"
)

type poolWrapper struct {
	s *Transport
}

func (w *poolWrapper) Close() {
}

type TransportPool struct  {
	p *pools.ResourcePool
}

func NewStoragePool(endpoint *url.URL, max int, timeout time.Duration) *TransportPool {
	newFn := func() (res pools.Resource, err error) {
		res = &poolWrapper{&Transport{endpoint}}
		return
	}
	return &TransportPool{pools.NewResourcePool(newFn, max, max, timeout)}
}

func (p *TransportPool) Take() (res *Transport, err error) {
	r, err := p.p.TryGet()
	if err != nil {
		return
	}
	res = r.(*poolWrapper).s
	return
}

func (p *TransportPool) Release(s *Transport) {
	p.p.Put(&poolWrapper{s})
}
