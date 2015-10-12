package server
import (
	"github.com/akaspin/bar/bard/storage"
	"github.com/akaspin/bar/proto"
	"net/rpc"
	"github.com/akaspin/bar/bard/service"
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
	service *rpc.Server
}

func NewBardServer(opts *BardServerOptions) (res *BardServer, err error) {
	res = &BardServer{BardServerOptions: opts}

	res.service = rpc.NewServer()
	rpcService := &service.Service{res.Info, res.StoragePool}
	if err = res.service.Register(rpcService); err != nil {
		return
	}

	return
}

func (s *BardServer) Start() (err error) {
	s.BardHttpServer = NewBardHttpServer(s.BardServerOptions, s.service)
	s.BardRPCServer = NewBardRPCServer(s.BardServerOptions, s.service)

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
