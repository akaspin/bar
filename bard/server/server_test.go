package server_test
import (
"testing"
"github.com/akaspin/bar/bard/storage"
	"time"
	"github.com/akaspin/bar/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/akaspin/bar/bard/server"
	"fmt"
"net/http"
	"github.com/akaspin/bar/proto"
)

func Test_StartStopBard(t *testing.T) {
	t.Skip()
	root := "test-start-stop"
	p := storage.NewStoragePool(
		storage.NewBlockStorageFactory(root, 2), 200, time.Minute)
	port, err := fixtures.GetOpenPort()
	assert.NoError(t, err)
	s := server.NewBardServer(&server.BardServerOptions{
		fmt.Sprintf(":%d", port),
		&proto.Info{},
		p,
	})

	go s.Start()
	time.Sleep(time.Second)

	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/v1/ping", port))
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	err = s.Stop()
	assert.NoError(t, err)
}
