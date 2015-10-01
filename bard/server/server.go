package server
import (
	"github.com/akaspin/bar/bard/storage"
	"net/http"
	"github.com/akaspin/bar/bard/handler"
	"github.com/tamtam-im/logx"
	"net"
)

type BardServerOptions struct  {
	Addr string
	ChunkSize int64
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

	mux := http.NewServeMux()

	mux.Handle("/v1/blob/upload/", &handler.SimpleUploadHandler{
		s.StoragePool, "/v1/blob/upload/"})
	mux.Handle("/v1/blob/check", &handler.CheckHandler{
		s.StoragePool})
	mux.Handle("/v1/tx/commit/declare", &handler.DeclareCommitTxHandler{
		s.StoragePool})
	mux.Handle("/v1/blob/download/", &handler.DownloadHandler{
		s.StoragePool, "/v1/blob/download/"})
	mux.Handle("/v1/ping", &handler.PingHandler{s.ChunkSize})

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
