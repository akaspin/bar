package fixtures

import (
	"fmt"
	"github.com/akaspin/bar/proto"
	"github.com/akaspin/bar/server"
	"github.com/akaspin/bar/server/front"
	"github.com/akaspin/bar/server/storage"
	"github.com/akaspin/bar/server/thrift"
	"github.com/tamtam-im/logx"
	"golang.org/x/net/context"
	"os"
	"path/filepath"
	"time"
)

type FixtureServer struct {
	server.Server
	*proto.ServerInfo
	Root string
}

func NewFixtureServer(name string) (res *FixtureServer, err error) {
	logx.SetLevel(logx.DEBUG)

	wd, _ := os.Getwd()
	rt := filepath.Join(wd, "testdata", name+"-srv")
	os.RemoveAll(rt)

	p := storage.NewBlockStorage(&storage.BlockStorageOptions{rt, 2, 32, 64})
	ports, err := GetOpenPorts(2)
	if err != nil {
		return
	}

	ctx := context.Background()
	info := &proto.ServerInfo{
		HTTPEndpoint: fmt.Sprintf("http://localhost:%d/v1", ports[0]),
		RPCEndpoints: []string{fmt.Sprintf("localhost:%d", ports[1])},
		ChunkSize:    1024 * 1024 * 2,
		PoolSize:     16,
		BufferSize:   1024 * 1024 * 8,
	}

	tServer := thrift.NewServer(ctx,
		&thrift.Options{info, fmt.Sprintf(":%d", ports[1])}, p)
	hServer := front.NewServer(ctx,
		&front.Options{info, fmt.Sprintf(":%d", ports[0]), ""}, p)

	res = &FixtureServer{
		server.NewCompositeServer(ctx, tServer, hServer),
		info,
		rt,
	}
	go res.Start()
	time.Sleep(time.Millisecond * 00)
	return
}

func (s *FixtureServer) Stop() {
	s.Server.Stop()
	os.RemoveAll(s.Root)
}

func ServerStoredName(root string, id string) string {
	wd, _ := os.Getwd()
	rt := filepath.Join(wd, "testdata", root+"srv")
	return filepath.Join(rt, "blobs", id[:2], id)
}
