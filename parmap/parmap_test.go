package parmap_test
import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
"math/rand"
	"time"
	"github.com/akaspin/bar/parmap"
)

func Test_ParMap1(t *testing.T)  {
	n := 1000

	req := map[string]interface{}{}
	for i := 0; i < n; i++ {
		req[fmt.Sprintf("%d", i)] = i
	}

	res, err := parmap.NewWorkerPool(128).RunBatch(
		parmap.Task{
			Map: req,
			Fn: func(k string, a interface{}) (res interface{}, err1 error) {
				return a, nil
			},
		})
	assert.NoError(t, err)
	assert.Len(t, res, 1000)
}

func Test_ParMap2(t *testing.T)  {
	n := 10

	req := map[string]interface{}{}
	for i := 0; i < n; i++ {
		req[fmt.Sprintf("%d", i)] = i
	}

	res, err := parmap.NewWorkerPool(128).RunBatch(parmap.Task{
		Map: req,
		Fn: func(k string, a interface{}) (res interface{}, err1 error) {
			return a, nil
		},
	})
	assert.NoError(t, err)
	assert.Len(t, res, 10)
}

func Test_ParMap_Nils(t *testing.T)  {
	n := 10

	req := map[string]interface{}{}
	for i := 0; i < n; i++ {
		req[fmt.Sprintf("%d", i)] = i
	}

	res, err := parmap.NewWorkerPool(128).RunBatch(parmap.Task{
		Map: req,
		Fn: func(k string, a interface{}) (res interface{}, err1 error) {
			return nil, nil
		},
	})
	assert.NoError(t, err)
	assert.Len(t, res, 10)
}




func Benchmark_ParMapWorkerPool(b *testing.B) {
	n := 1000000
	var v string
	req := map[string]interface{}{}
	for i := 0; i < n; i++ {
		v = fmt.Sprintf("%d", i)
		req[v] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		res, err := parmap.NewWorkerPool(128).RunBatch(parmap.Task{
			Map: req,
			Fn: func(k string, a interface{}) (res interface{}, err1 error) {
				return op(a.(int)), nil
			},
		})
		b.StopTimer()
		assert.NoError(b, err)
		assert.Len(b, res, n)
	}
}

func Benchmark_ParMapWG(b *testing.B) {
	n := 1000000
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		res := parmapWG(n)
		b.StopTimer()
		assert.NoError(b, checkRes(res, n))
	}
}

func Benchmark_ParMapChan(b *testing.B) {
	n := 1000000
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		res := parmapChan(n)
		b.StopTimer()
		assert.NoError(b, checkRes(res, n))
	}
}

func parmapChan(n int) (res map[string]string) {
	res = map[string]string{}
	resChan := make(chan string, 8)

	for i := 0; i < n; i++ {
		go func(i int) {
			resChan <- op(i)
		}(i)
	}

	for i := 0; i < n; i++ {
		r :=<- resChan
		res[r] = r
		//		if len(res) == n {
		//			break
		//		}
	}
	return
}

func parmapWG(n int) (res map[string]string) {
	res = map[string]string{}

	var wg sync.WaitGroup
	var mu sync.RWMutex

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			r := op(i)
			mu.Lock()
			res[r] = r
			mu.Unlock()
		}(i)
	}
	wg.Wait()
	return
}

func op(i int) string {
	rand.Seed(time.Now().Unix())
	time.Sleep(time.Microsecond * time.Duration(rand.Int31n(1000)))
	return fmt.Sprintf("%d", i)
}

func checkRes(res map[string]string, n int) (err error) {
	for i := 0; i < n; i++ {
		key := fmt.Sprintf("%d", i)
		r, ok := res[key]
		if !ok || r != key {
			err = fmt.Errorf("bad i: %d", i)
			return
		}
	}
	return
}