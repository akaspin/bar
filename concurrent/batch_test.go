package concurrent_test
import (
	"testing"
	"math/rand"
	"time"
	"fmt"
	"github.com/akaspin/bar/concurrent"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func Test_Concurrent_Pool1(t *testing.T) {
	n := 100

	var req, res []interface{}
	for i := 0; i < n; i++ {
		req = append(req, i)
	}

	pool := concurrent.NewPool(16)
	defer pool.Close()

	err := pool.Do(
		func(ctx context.Context, in interface{}) (out interface{}, err error) {
			out = fmt.Sprintf("%d", in.(int))
			return
		},
		&req, &res,
		concurrent.DefaultBatchOptions(),
	)
	assert.NoError(t, err)
	assert.Len(t, res, 100)
}

func Test_Concurrent_Pool_Error(t *testing.T) {
	n := 100

	var req []interface{}
	var res []interface{}
	for i := 0; i < n; i++ {
		req = append(req, i)
	}

	pool := concurrent.NewPool(16)
	defer pool.Close()
	err := pool.Do(
		func(ctx context.Context, in interface{}) (out interface{}, err error) {
			if in.(int) >= 0 && in.(int) < 30 {
				err = fmt.Errorf("test error")
				return
			}
			out = fmt.Sprintf("%d", in.(int))
			return
		},
		&req,
		&res,
		concurrent.DefaultBatchOptions(),
	)
	assert.Error(t, err)
	assert.NotEqual(t, len(res), 100)
}

func Test_Concurrent_Pool_IgnoreErrors(t *testing.T) {
	n := 100

	var req []interface{}
	var res []interface{}
	for i := 0; i < n; i++ {
		req = append(req, i)
	}

	pool := concurrent.NewPool(16)
	defer pool.Close()
	err := pool.Do(
		func(ctx context.Context, in interface{}) (out interface{}, err error) {
			if in.(int) >= 0 && in.(int) < 30 {
				err = fmt.Errorf("test error")
				return
			}
			out = fmt.Sprintf("%d", in.(int))
			return
		},
		&req,
		&res,
		concurrent.DefaultBatchOptions().AllowErrors(),
	)
	assert.Error(t, err)
	assert.Equal(t, len(res), 70)
}

func Test_Concurrent_Pool_Hangs(t *testing.T) {
	n := 100

	var req []interface{}
	var res []interface{}
	for i := 0; i < n; i++ {
		req = append(req, i)
	}

	pool := concurrent.NewPool(16)
	defer pool.Close()
	time.Sleep(time.Second)
	assert.Equal(t, 0, pool.Busy())

	err := pool.Do(
		func(ctx context.Context, in interface{}) (out interface{}, err error) {
			if in.(int) == 90 {
				time.Sleep(time.Second * 5)
				return "ok", nil
			}
			return nil, fmt.Errorf("break it")
		},
		&req, &res, concurrent.DefaultBatchOptions(),
	)

	time.Sleep(time.Second)
	assert.Equal(t, 0, pool.Busy())
	assert.Error(t, err)
	assert.Equal(t, len(res), 0)
}

func op(i int) string {
	rand.Seed(time.Now().Unix())
	time.Sleep(time.Microsecond * time.Duration(rand.Int31n(1000)))
	return fmt.Sprintf("%d", i)
}