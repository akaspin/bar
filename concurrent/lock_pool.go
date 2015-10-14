package concurrent
import (
	"github.com/vireshas/minimal_vitess_pool/pools"
	"time"
)

type Lock struct {
	pool *LocksPool
	IsClosed bool
}

func (l *Lock) Close() {
	if !l.IsClosed {
		l.IsClosed = true
		l.pool.pool.Put(nil)
	}
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
//	go func() {
//		select {
//		case <-time.After(p.timeout):
//			res.Close()
//		}
//	}()

	return
}

func (p *LocksPool) factory() (res pools.Resource, err error) {
	res = &Lock{p, false}
	return
}