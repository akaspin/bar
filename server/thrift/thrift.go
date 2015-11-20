package thrift

import (
	"github.com/akaspin/bar/proto/wire"
	"github.com/akaspin/bar/server/storage"
	t_thrift "github.com/apache/thrift/lib/go/thrift"
	"github.com/tamtam-im/logx"
	"golang.org/x/net/context"
)

type ThriftServer struct {
	ctx     context.Context
	cancel  context.CancelFunc
	options *Options
	*Handler
	t_thrift.TServer
}

func NewServer(ctx context.Context, opts *Options, stor storage.Storage) (res *ThriftServer) {
	res = &ThriftServer{
		ctx:     ctx,
		options: opts,
	}
	res.ctx, res.cancel = context.WithCancel(ctx)
	res.Handler = NewHandler(ctx, opts.Info, stor)

	return
}

func (s *ThriftServer) Start() (err error) {
	processor := wire.NewBarProcessor(s.Handler)

	var transport t_thrift.TServerTransport
	if transport, err = t_thrift.NewTServerSocket(s.options.Bind); err != nil {
		return
	}
	protoFactory := t_thrift.NewTBinaryProtocolFactoryDefault()
	transportFactory := t_thrift.NewTBufferedTransportFactory(
		s.BufferSize)
	s.TServer = t_thrift.NewTSimpleServer4(processor, transport,
		transportFactory, protoFactory)

	logx.Debugf("thrift listening at %s", s.options.Info.RPCEndpoints[0])

	errChan := make(chan error, 1)
	go func() {
		errChan <- s.TServer.Serve()

	}()
	defer s.TServer.Stop()
	defer s.cancel()

	select {
	case <-s.ctx.Done():
		return
	case err = <-errChan:
		return
	}
}

func (s *ThriftServer) Stop() (err error) {
	s.cancel()
	logx.Debugf("thrift stopped at %s", s.RPCEndpoints[0])
	return
}
