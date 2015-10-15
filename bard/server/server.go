package server
import (
	"github.com/akaspin/bar/bard/storage"
	"github.com/akaspin/bar/proto"
)


type BardServerOptions struct  {
	// Client info
	*proto.ServerInfo
	HttpBind string
	RPCBind string
	storage.Storage
	BarExe string
}

type BardServer struct {
	*BardServerOptions
	*BardHttpServer
	*ThriftServer
}

func NewBardServer(opts *BardServerOptions) (res *BardServer) {
	res = &BardServer{BardServerOptions: opts}
	return
}

func (s *BardServer) Start() (err error) {
	s.BardHttpServer = NewBardHttpServer(s.BardServerOptions)
	s.ThriftServer =NewThriftServer(s.BardServerOptions)

	errChan := make(chan error, 1)

	go func() {
		errChan <- s.ThriftServer.Start()
	}()

	go func() {
		errChan <- s.BardHttpServer.Start()
	}()

	err = <- errChan

	return
}

func (s *BardServer) Stop() (err error) {
	s.BardHttpServer.Stop()
	s.ThriftServer.Stop()
	return
}
