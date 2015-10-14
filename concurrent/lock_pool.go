package concurrent
import (
	"github.com/vireshas/minimal_vitess_pool/pools"
	"time"
)

type Lock struct {
	pool *LocksPool
}

func (l *Lock) Close() {}

func (l *Lock) Release() {
	l.pool.pool.Put(l)
}

type LocksPool struct {
	pool *pools.ResourcePool
	timeout time.Duration
}

func NewLockPool(n int, timeout time.Duration) (res *LocksPool) {
	res = &LocksPool{
		timeout: timeout,
	}
	res.pool = pools.NewResourcePool(res.factory, n, n, timeout)
	return 
}

func (p *LocksPool) Take() (res *Lock, err error) {
	r, err := p.pool.Get(p.timeout)
	if err != nil {
		return
	}
	res = r.(*Lock)
	return
}

func (p *LocksPool) factory() (res pools.Resource, err error) {
	res = &Lock{p}
	return
}