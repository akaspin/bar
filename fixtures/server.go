package fixtures
import (
	"testing"
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

func RunServer(t *testing.T, root string) (httpEndpoint string, rpcEndpoints string, stop func())  {
	logx.SetLevel(logx.DEBUG)

	wd, _ := os.Getwd()
	rt := filepath.Join(wd, "testdata", root + "srv")
	os.RemoveAll(rt)

	p := storage.NewBlockStorage(&storage.BlockStorageOptions{rt, 2, 32, 64})
	ports, err := GetOpenPorts(2)
	assert.NoError(t, err)

	httpEndpoint = fmt.Sprintf("http://127.0.0.1:%d/v1", ports[0])
	rpcEndpoints = fmt.Sprintf("localhost:%d", ports[1])
	srv, err := server.NewBardServer(&server.BardServerOptions{
		HttpBind: fmt.Sprintf(":%d", ports[0]),
		RPCBind: fmt.Sprintf(":%d", ports[1]),
		Info: &proto.ServerInfo{
			HTTPEndpoint: fmt.Sprintf("http://localhost:%d/v1", ports[0]),
			RPCEndpoints: []string{fmt.Sprintf("localhost:%d", ports[1])},
			ChunkSize: 1024 * 1024 * 2,
			PoolSize: 16,
			BufferSize: 1024 * 1024 * 8,
		},
		Storage: p,
		BarExe: "",
	})
	if err != nil {
		return
	}

	go srv.Start()
	time.Sleep(time.Millisecond * 300)
	assert.NoError(t, err)
	stop = func() {
		srv.Stop()
		os.RemoveAll(rt)
	}
	return
}

func ServerStoredName(root string, id string) string {
	wd, _ := os.Getwd()
	rt := filepath.Join(wd, "testdata", root + "srv")
	return filepath.Join(rt, "blobs", id[:2], id)
}
