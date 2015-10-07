package server
import (
	"github.com/akaspin/bar/bard/storage"
	"net/http"
	"github.com/tamtam-im/logx"
	"net"
	"github.com/akaspin/bar/bard/service"
	"net/rpc"
	"github.com/akaspin/bar/proto"
)

type BardServerOptions struct  {
	Addr string
	Info *proto.Info
	StoragePool *storage.StoragePool
}

type BardServer struct {
	*BardServerOptions
	l net.Listener
}

func NewBardServer(opts *BardServerOptions) *BardServer {
	return &BardServer{BardServerOptions: opts}
}

func (s *BardServer) Start() (err error) {
	s.l, err = net.Listen("tcp", s.Addr)
	if err != nil {
		return
	}

//	mux.Handle("/v1/blob/download/", &handler.DownloadHandler{
//		s.StoragePool, "/v1/blob/download/"})

	rpcSvr := rpc.NewServer()
	rpcService := &service.Service{s.Info, s.StoragePool}
	if err = rpcSvr.Register(rpcService); err != nil {
		return
	}

	mux := http.NewServeMux()
	mux.Handle("/v1/rpc", rpcSvr)

	logx.Debugf("bard serving at http://%s/v1", s.Addr)
	srv := &http.Server{Handler:mux}
	err = srv.Serve(s.l)
	return
}

func (s *BardServer) Stop() (err error) {
	err = s.l.Close()
	if err != nil {
		return
	}
	logx.Debugf("http://%s/v1 bye!", s.Addr)
	return
}
