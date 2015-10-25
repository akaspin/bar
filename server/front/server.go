package front

import (
	"net"
	"net/http"
	"github.com/tamtam-im/logx"
	"golang.org/x/net/context"
	"github.com/akaspin/bar/server/storage"
)

type Server struct  {
	ctx context.Context
	cancel context.CancelFunc

	options *Options
	storage storage.Storage

	net.Listener
}

func NewServer(ctx context.Context, opts *Options, stor storage.Storage) (res *Server) {
	res = &Server{
		options: opts,
		storage: stor,
	}
	res.ctx, res.cancel = context.WithCancel(ctx)
	return
}

func (s *Server) Start() (err error) {
	if s.Listener, err = net.Listen("tcp", s.options.Bind); err != nil {
		return
	}

	hs, err := NewHandlers(s.ctx, s.options, s.storage)
	if err != nil {
		return
	}

	// make http frontend
	mux := http.NewServeMux()
	mux.HandleFunc("/", hs.HandleFront)
	mux.HandleFunc("/v1/win/bar-export.bat", hs.HandleExportBat)
	mux.HandleFunc("/v1/win/bar-import/", hs.HandleImportBat)
	mux.HandleFunc("/v1/win/bar.exe", hs.HandleBarExe)
	mux.HandleFunc("/v1/spec/", hs.HandleSpec)
	logx.Debugf("bard http serving at http://%s/v1", s.options.Bind)
	srv := &http.Server{Handler:mux}

	errChan := make(chan error, 1)
	go func() {
		errChan <- srv.Serve(s.Listener)

	}()
	defer s.Listener.Close()
	defer s.cancel()

	select {
	case <-s.ctx.Done():
		break
	case err = <- errChan:
		break
	}
	return
}

func (s *Server) Stop() (err error) {
	s.cancel()
	logx.Debugf("http %s closed", s.options.Bind)
	return
}

