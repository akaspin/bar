package concurrent
import (
	"github.com/vireshas/minimal_vitess_pool/pools"
	"time"
	"sync/atomic"
)

type dummy struct {}

func (d *dummy) Close() {}

type Lock struct {
	pool *LocksPool
	n *int32
}

func newLock(p *LocksPool, n int) (res *Lock, err error) {
	var n1 int32
	res = &Lock{p, &n1}
	for i := 0; i < n; i++ {
		atomic.AddInt32(res.n, 1)
		if _, err = p.pool.Get(p.timeout); err != nil {
			res.close()
			return
		}
	}
	return
}

func (l *Lock) Close() {
	l.close()
}

func (l *Lock) close() {
	n1 := atomic.LoadInt32(l.n)
	var i int32
	for i = 0; i < n1; i++ {
		l.pool.pool.Put(nil)
		atomic.AddInt32(l.n, -1)
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

func (p *LocksPool) Available() (res int64) {
	return p.pool.Available()
}

func (p *LocksPool) With(n int, fn func()) (err error) {
	l, err := p.TakeN(n)
	if err != nil {
		return
	}
	defer l.Close()
	fn()
	return
}

func (p *LocksPool) TakeN(n int) (res *Lock, err error)  {
	return newLock(p, n)
}

// Take single lock
func (p *LocksPool) Take() (res *Lock, err error) {
	res, err = newLock(p, 1)
	return
}

func (p *LocksPool) Close() {
	p.pool.Close()
}

func (p *LocksPool) factory() (res pools.Resource, err error) {
	res = &dummy{}
	return
}