package concurrent_test
import (
	"testing"
	"golang.org/x/net/context"
	"github.com/stretchr/testify/assert"
)

func Test_Concurrent_Cancel(t *testing.T) {
	baseCtx, baseCancel := context.WithCancel(context.Background())

	resChan := make(chan int, 100)

	for i := 0; i < 100; i++ {
	 	go func(ctx context.Context, i int) {
		    for {
			    select {
			    case <-ctx.Done():
					resChan <- i
					return
			    }
		    }
	    }(baseCtx, i)
	}

	res := map[int]struct{}{}
	go baseCancel()
	for i := 0; i < 100; i++ {
		select {
		case ii := <- resChan:
			res[ii] = struct{}{}
		}
	}
	assert.Len(t, res, 100)
}