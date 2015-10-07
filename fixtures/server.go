package fixtures
import (
	"testing"
	"net/url"
	"github.com/akaspin/bar/bard/storage"
	"time"
	"github.com/stretchr/testify/assert"
	"github.com/akaspin/bar/bard/server"
	"fmt"
	"path/filepath"
	"github.com/tamtam-im/logx"
	"github.com/akaspin/bar/proto"
)

func RunServer(t *testing.T, root string) (endpoint *url.URL, stop func() error)  {
	logx.SetLevel(logx.DEBUG)

	p := storage.NewStoragePool(storage.NewBlockStorageFactory(root, 2), 200, time.Minute)
	port, err := GetOpenPort()
	assert.NoError(t, err)

	endpoint, err = url.Parse(fmt.Sprintf("http://127.0.0.1:%d/v1", port))
	srv := server.NewBardServer(&server.BardServerOptions{
		fmt.Sprintf(":%d", port),
		&proto.Info{
			[]string{endpoint.String()},
			1024 * 1024 * 2, 16}, p,
	})

	go srv.Start()
	time.Sleep(time.Millisecond * 200)
	assert.NoError(t, err)
	stop = srv.Stop
	return
}

func ServerStoredName(root string, id string) string {
	return filepath.Join(root, "blobs", id[:2], id)
}
