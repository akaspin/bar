package fixtures_test
import (
	"testing"
	"sync"
	"fmt"
	"math/rand"
	"time"
	"github.com/stretchr/testify/assert"
)



func Benchmark_ParMapWG(b *testing.B) {
	n := 1000000
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		res := parmapWG(n)
		b.StopTimer()
		assert.NoError(b, checkRes(res, n))
	}
}

func Benchmark_BarMapChan(b *testing.B) {
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
	resChan := make(chan string, 1)

	for i := 0; i < n; i++ {
		go func(i int) {
			resChan <- op(i)
		}(i)
	}

	for {
		r :=<- resChan
		res[r] = r
		if len(res) == n {
			break
		}
	}
	return
}

func parmapWG(n int) (res map[string]string) {
	res = map[string]string{}

	var wg sync.WaitGroup
	var mu sync.RWMutex

	for i := 0; i < n; i++ {
		wg.Add(1)
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