package server
import (
	"github.com/akaspin/bar/bard/storage"
	"github.com/akaspin/bar/proto"
)

type BardServerOptions struct  {
	HttpAddr string
	RPCAddr string
	Info *proto.Info
	StoragePool *storage.StoragePool
	BarExe string
}

type BardServer struct {
	*BardServerOptions
	*BardHttpServer
	*BardRPCServer
}

func NewBardServer(opts *BardServerOptions) *BardServer {
	return &BardServer{BardServerOptions: opts}
}

func (s *BardServer) Start() (err error) {
	s.BardHttpServer = NewBardHttpServer(s.BardServerOptions)
	s.BardRPCServer = NewBardRPCServer(s.BardServerOptions)

	errChan := make(chan error, 1)

	go func() {
		errChan <- s.BardRPCServer.Start()
	}()

	go func() {
		errChan <- s.BardHttpServer.Start()
	}()

	err = <- errChan

	return
}

func (s *BardServer) Stop() (err error) {
	s.BardHttpServer.Stop()
	s.BardRPCServer.Stop()
	return
}
