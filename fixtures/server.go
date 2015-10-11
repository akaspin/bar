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
	ports, err := GetOpenPorts(2)
	assert.NoError(t, err)

	endpoint, _ = url.Parse(fmt.Sprintf("http://127.0.0.1:%d/v1", ports[0]))
	srv := server.NewBardServer(&server.BardServerOptions{
		HttpAddr: fmt.Sprintf(":%d", ports[1]),
		RPCAddr: fmt.Sprintf(":%d", ports[0]),
		Info: &proto.Info{
			HTTPEndpoint: fmt.Sprintf("http://localhost:%d/v1", ports[1]),
			Endpoint: fmt.Sprintf("http://localhost:%d/v1", ports[0]),
			ChunkSize: 1024 * 1024 * 2,
			PoolSize: 16,
		},
		StoragePool: p,
		BarExe: "",
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
