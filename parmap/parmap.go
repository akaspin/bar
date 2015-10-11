package parmap
import (
	"golang.org/x/net/context"
	"fmt"
)

type taskRequest struct  {
	ResChan chan *taskResult
	Key string
	Fn func(key string, arg interface{}) (res interface{}, err error)
	Args interface{}
}

type taskResult struct {
	Key string
	Value interface{}
	Error error
}

func serve(
	ctx context.Context,
	cancel context.CancelFunc,
	taskChan <-chan *taskRequest,
) {
	for {
		select {
		case task := <- taskChan:
			res, err := task.Fn(task.Key, task.Args)
			task.ResChan <- &taskResult{task.Key, res, err}
		case <-ctx.Done():
			return
		}
	}
}

type Task struct {

	// Data
	Map map[string]interface{}

	// Mapping fn
	Fn func(string, interface{}) (interface{}, error)

	// Do not fail on errors
	IgnoreErrors bool
}

type ParMap struct {
	taskChan chan *taskRequest
	ctx context.Context
	cancel context.CancelFunc
	size int
}

func NewWorkerPool(size int) (res *ParMap)  {
	res = &ParMap{
		taskChan: make(chan *taskRequest, size),
		size: size,
	}
	res.ctx, res.cancel = context.WithCancel(context.Background())
	for i := 0; i <= size; i++ {
		wCtx, wCancel := context.WithCancel(res.ctx)
		go serve(wCtx, wCancel, res.taskChan,)
	}
	return
}


func (p ParMap) RunOne(
	key string,
	arg interface{},
	fn func(key string, arg interface{}) (res interface{}, err error),
) (res interface{}, err error) {
	resChan := make(chan *taskResult)
	ctx, _ := context.WithCancel(p.ctx)

	go func(key string, arg interface{}) {
		p.taskChan <- &taskRequest{resChan, key, fn, arg}
	}(key, arg)

	select {
	case <-ctx.Done():
		break
	case r := <-resChan:
		res, err = r.Value, r.Error
	}

	return
}

func (p *ParMap) RunBatch(task Task) (res map[string]interface{}, err error) {
	res = map[string]interface{}{}
	var errs []error
	resChan := make(chan *taskResult)
	ctx, _ := context.WithCancel(p.ctx)

	for k, a := range task.Map {
		go func(key string, arg interface{}) {
			p.taskChan <- &taskRequest{resChan, key, task.Fn, arg}
		}(k, a)
	}

	for i := 0; i < len(task.Map); i++ {
		select {
		case <-ctx.Done():
			break
		case r := <-resChan:
			if r.Error != nil {
				errs = append(errs, r.Error)
				if !task.IgnoreErrors {
					break
				}
			} else {
				res[r.Key] = r.Value
			}
		}
	}

	if len(errs) > 0 {
		err = fmt.Errorf("%s", errs)
	}
	return
}

func (p *ParMap) Close() {
	p.cancel()
}