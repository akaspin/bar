package server

import "golang.org/x/net/context"

type Server interface {
	Start() (err error)
	Stop() (err error)
}

type CompositeServer struct {
	ctx     context.Context
	cancel  context.CancelFunc
	servers []Server
}

func NewCompositeServer(ctx context.Context, servers ...Server) (res *CompositeServer) {
	res = &CompositeServer{servers: servers}
	res.ctx, res.cancel = context.WithCancel(ctx)
	return
}

func (s *CompositeServer) Start() (err error) {
	errChan := make(chan error, len(s.servers))

	for _, srv := range s.servers {
		go func(srv Server) {
			errChan <- srv.Start()
		}(srv)
	}
	defer s.cancel()

loop:
	for {
		select {
		case <-s.ctx.Done():
			break loop
		case err = <-errChan:
			s.cancel()
			break loop
		}
	}

	return
}

func (s *CompositeServer) Stop() (err error) {
	s.cancel()
	return
}
