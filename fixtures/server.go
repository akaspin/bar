package fixtures
import (
	"github.com/akaspin/bar/bard/storage"
	"time"
	"github.com/akaspin/bar/bard/server"
	"fmt"
	"path/filepath"
	"github.com/tamtam-im/logx"
	"github.com/akaspin/bar/proto"
	"os"
)

type FixtureServer struct {
	*server.BardServer
	Root string
}

func NewFixtureServer(name string) (res *FixtureServer, err error) {
	logx.SetLevel(logx.DEBUG)

	wd, _ := os.Getwd()
	rt := filepath.Join(wd, "testdata", name + "-srv")
	os.RemoveAll(rt)

	p := storage.NewBlockStorage(&storage.BlockStorageOptions{rt, 2, 32, 64})
	ports, err := GetOpenPorts(2)
	if err != nil {
		return
	}

	res = &FixtureServer{
		BardServer: server.NewBardServer(&server.BardServerOptions{
			HttpBind: fmt.Sprintf(":%d", ports[0]),
			RPCBind: fmt.Sprintf(":%d", ports[1]),
			ServerInfo: &proto.ServerInfo{
				HTTPEndpoint: fmt.Sprintf("http://localhost:%d/v1", ports[0]),
				RPCEndpoints: []string{fmt.Sprintf("localhost:%d", ports[1])},
				ChunkSize: 1024 * 1024 * 2,
				PoolSize: 16,
				BufferSize: 1024 * 1024 * 8,
			},
			Storage: p,
			BarExe: "",
		}),
		Root: rt,
	}
	go res.BardServer.Start()
	time.Sleep(time.Millisecond * 200)
	return
}

func (s *FixtureServer) Stop() {
	s.BardServer.Stop()
	os.RemoveAll(s.Root)
}

func ServerStoredName(root string, id string) string {
	wd, _ := os.Getwd()
	rt := filepath.Join(wd, "testdata", root + "srv")
	return filepath.Join(rt, "blobs", id[:2], id)
}
