package concurrent
import (
	"fmt"
	"time"
	"sync/atomic"
)

type jobRequest struct {
	fn func(interface{}) (interface{}, error)
	arg interface{}
	resChan chan<- interface{}
	errChan chan<- error
	cancelChan chan struct{}
}

type worker struct {
	jobChan chan *jobRequest
	stopChan chan struct{}
	timeout time.Duration
	busy *int32
}

func newWorker(jobChan chan *jobRequest, timeout time.Duration, busy *int32) (res *worker) {
	res = &worker{
		jobChan,
		make(chan struct{}, 1),
		timeout,
		busy,
	}
	go res.work()
	return
}

func (w *worker) work() {
	for {
		select {
		case job := <-w.jobChan:
			atomic.AddInt32(w.busy, 1)
			if w.do(job) {
				return
			}
		case <- w.stopChan:
			return
		}
		atomic.AddInt32(w.busy, -1)
	}
}

func (w *worker) do(job *jobRequest) (stop bool) {
	// if fn hangs - just leave this channels

	resChan := make(chan interface{}, 1)
	errChan := make(chan error, 1)
	go func() {
		if res, err := job.fn(job.arg); err != nil {
			errChan <- err
		} else {
			resChan <- res
		}
	}()
	select {
	case <- time.After(w.timeout):
		job.errChan <- fmt.Errorf("timeout %v", job.arg)
		return
	case <- job.cancelChan:
		return
	case res := <- resChan:
		job.resChan <- res
		return
	case err := <- errChan:
		job.errChan <- err
		return
	case <- w.stopChan:
		stop = true
		return
	}
	return
}

func (w *worker) close() {
	go func() {
		w.stopChan <- struct{}{}
	}()
}

type Pool struct {
	workers []*worker
	jobChan chan *jobRequest
	busy int32
}

func NewPool(n int) (res *Pool) {
	res = &Pool{
		jobChan: make(chan *jobRequest, n),
		busy: 0,
	}
	for i := 0; i < n; i++ {
		res.workers = append(res.workers, newWorker(res.jobChan, time.Hour, &res.busy))
	}
	return
}

func (p *Pool) Do(
	fn func(interface{}) (interface{}, error),
	in *[]interface{},
	out *[]interface{},
	ignoreErrors bool,
	acceptNils bool,
) (err error) {
	n := len(*in)
	resChan := make(chan interface{}, n)
	errChan := make(chan error, n)
	var errs []error

	cancels := make([]chan struct{}, n)

	for _, arg := range *in {
		cancelChan := make(chan struct{}, 1)
		cancels = append(cancels, cancelChan)
		go func(arg interface{}) {
			p.jobChan <- &jobRequest{fn, arg, resChan, errChan, cancelChan}
		}(arg)
	}

	poll:
	for i := 0; i < n; i++ {
		select {
		case err1 := <- errChan:
			if !ignoreErrors {
				err = err1
				break poll
			}
			errs = append(errs, err1)
		case res := <- resChan:
			if res != interface{}(nil) || acceptNils {
				*out = append(*out, res)
			}
		}
	}
	for _, cancelChan := range cancels {
		go func(c chan struct{}) {
			c <- struct{}{}
		}(cancelChan)
	}

	if len(errs) > 0 {
		err = fmt.Errorf("%s", errs)
	}
	return
}

func (p *Pool) Busy() int {
	return int(atomic.LoadInt32(&p.busy))

}

func (p *Pool) Close() {
	for _, w := range p.workers {
		w.close()
	}
}