package concurrent_test
import (
	"testing"
	"github.com/akaspin/bar/concurrent"
	"time"
	"github.com/stretchr/testify/assert"
)


func Test_Concurrent_LocksPool1(t *testing.T)  {
	p := concurrent.NewLockPool(1000, time.Minute)
	// take 600 locks
	for i := 100; i < 700; i += 100 {
		func(i int) {
			l, err := p.TakeN(100)
			assert.NoError(t, err)
			defer l.Close()
		}(i)
	}
	assert.EqualValues(t, 1000, p.Available())
}

func Test_Concurrent_LocksPool_With(t *testing.T)  {
	p := concurrent.NewLockPool(1000, time.Minute)
	for i := 100; i < 700; i += 100 {
		assert.NoError(t, p.With(100, func() {  }))
	}
	assert.EqualValues(t, 1000, p.Available())
}