package server
import (
	"net"
	"net/rpc"
	"github.com/tamtam-im/logx"
	"net/http"
)

type BardRPCServer struct  {
	*BardServerOptions
	service *rpc.Server
	net.Listener
}

func NewBardRPCServer(opts *BardServerOptions, service *rpc.Server) *BardRPCServer {
	return &BardRPCServer{BardServerOptions: opts, service: service}
}

func (s *BardRPCServer) Start() (err error) {
	if s.Listener, err = net.Listen("tcp", s.RPCAddr); err != nil {
		return
	}

	logx.Debugf("bard RPC serving at tcp://%s", s.RPCAddr)
	mux := http.NewServeMux()
	mux.Handle("/v1/rpc", s.service)
	logx.Debugf("bard http serving at http://%s/v1", s.HttpAddr)
	srv := &http.Server{Handler:mux}
	err = srv.Serve(s.Listener)
	return
}

func (s *BardRPCServer) Stop() (err error) {
	logx.Tracef("closing rpc tcp://%s", s.HttpAddr)
	if err = s.Listener.Close(); err != nil {
		return
	}
	logx.Debugf("rpc tcp://%s is closed", s.HttpAddr)
	return
}