package server
import (
	"github.com/akaspin/bar/bard/storage"
	"github.com/akaspin/bar/proto"
	"net/rpc"
	"github.com/akaspin/bar/bard/service"
)


type BardServerOptions struct  {
	// Client info
	Info *proto.Info
	HttpBind string
	RPCBind string
	storage.Storage
	BarExe string
}

type BardServer struct {
	*BardServerOptions
	*BardHttpServer
	service *rpc.Server
}

func NewBardServer(opts *BardServerOptions) (res *BardServer, err error) {
	res = &BardServer{BardServerOptions: opts}

	res.service = rpc.NewServer()
	rpcService := &service.Service{res.Info, res.Storage}
	if err = res.service.Register(rpcService); err != nil {
		return
	}

	return
}

func (s *BardServer) Start() (err error) {
	s.BardHttpServer = NewBardHttpServer(s.BardServerOptions, s.service)

	errChan := make(chan error, 1)

	go func() {
		errChan <- s.BardHttpServer.Start()
	}()

	err = <- errChan

	return
}

func (s *BardServer) Stop() (err error) {
	s.BardHttpServer.Stop()
	return
}
