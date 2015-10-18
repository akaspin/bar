package server
import (
	"net"
	"net/http"
	"github.com/akaspin/bar/bard/handler"
	"github.com/tamtam-im/logx"
)

type BardHttpServer struct  {
	*BardServerOptions
	net.Listener
}

func NewBardHttpServer(opts *BardServerOptions) *BardHttpServer {
	return &BardHttpServer{BardServerOptions: opts}
}

func (s *BardHttpServer) Start() (err error) {
	s.Listener, err = net.Listen("tcp", s.HttpBind)
	if err != nil {
		return
	}

	// make http frontend
	mux := http.NewServeMux()
	mux.Handle("/", &handler.FrontHandler{s.ServerInfo})
	mux.Handle("/v1/win/bar-export.bat", &handler.ExportBatHandler{s.ServerInfo})
	mux.Handle("/v1/win/bar-import/", &handler.ImportBatHandler{s.ServerInfo})
	mux.Handle("/v1/win/bar.exe", &handler.ExeHandler{s.BarExe})
	mux.Handle("/v1/spec/", &handler.SpecHandler{
		s.Storage, s.ServerInfo, s.BarExe})
//	mux.Handle("/v1/rpc", s.service)
	logx.Debugf("bard http serving at http://%s/v1", s.HttpBind)
	srv := &http.Server{Handler:mux}
	err = srv.Serve(s.Listener)
	return
}

func (s *BardHttpServer) Stop() (err error) {
	logx.Tracef("closing http://%s/v1", s.HttpBind)
	if err = s.Listener.Close(); err != nil {
		return
	}
	logx.Debugf("http %s closed", s.HttpBind)
	return
}