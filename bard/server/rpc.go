package server
import (
	"net"
	"net/rpc"
	"github.com/akaspin/bar/bard/service"
"net/http"
	"github.com/tamtam-im/logx"
)

type BardRPCServer struct  {
	*BardServerOptions
	net.Listener
}

func NewBardRPCServer(opts *BardServerOptions) *BardRPCServer {
	return &BardRPCServer{BardServerOptions: opts}
}

func (s *BardRPCServer) Start() (err error) {
	if s.Listener, err = net.Listen("tcp", s.RPCAddr); err != nil {
		return
	}

	// make rpc server
	rpcSvr := rpc.NewServer()
	rpcService := &service.Service{s.Info, s.StoragePool}
	if err = rpcSvr.Register(rpcService); err != nil {
		return
	}

	// make http frontend
	mux := http.NewServeMux()
	mux.Handle("/v1/rpc", rpcSvr)

	logx.Debugf("bard RPC serving at http://%s/v1", s.RPCAddr)
	srv := &http.Server{Handler:mux}
	err = srv.Serve(s.Listener)
	return
}

func (s *BardRPCServer) Stop() (err error) {
	logx.Tracef("closing http://%s/v1", s.HttpAddr)
	if err = s.Listener.Close(); err != nil {
		return
	}
	logx.Debugf("http %s closed", s.HttpAddr)
	return
}