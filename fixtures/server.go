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
	"os"
)

func RunServer(t *testing.T, root string) (endpoint *url.URL, stop func() error)  {
	logx.SetLevel(logx.DEBUG)

	wd, _ := os.Getwd()
	rt := filepath.Join(wd, "testdata", root + "srv")
	os.RemoveAll(rt)

	p := storage.NewStoragePool(storage.NewBlockStorageFactory(rt, 2), 200, time.Minute)
	port, err := GetOpenPort()
	assert.NoError(t, err)

	endpoint, err = url.Parse(fmt.Sprintf("http://127.0.0.1:%d/v1", port))
	srv := server.NewBardServer(&server.BardServerOptions{
		fmt.Sprintf(":%d", port),
		&proto.Info{
			fmt.Sprintf("http://localhost:%d/v1"),
			fmt.Sprintf("http://localhost:%d/v1"),
			1024 * 1024 * 2, 16},
		p,
		"",
	})

	go srv.Start()
	time.Sleep(time.Millisecond * 300)
	assert.NoError(t, err)
	stop = srv.Stop
	return
}

func ServerStoredName(root string, id string) string {
	wd, _ := os.Getwd()
	rt := filepath.Join(wd, "testdata", root + "srv")
	return filepath.Join(rt, "blobs", id[:2], id)
}
