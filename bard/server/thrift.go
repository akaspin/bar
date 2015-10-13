package server

import (
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/akaspin/bar/proto/bar"
)

type ThriftServer struct  {
	*BardServerOptions
	Server thrift.TServer
}

func (s *ThriftServer) Start() {
	processor := bar.NewBarProcessor()
	go s.start()
}

func (s *ThriftServer) start(processor thrift.TProcessor) (err error) {
	var transport thrift.TServerTransport

	if transport, err = thrift.NewTServerSocket(s.RPCBind); err != nil {
		return
	}
	protoFactory := thrift.NewTBinaryProtocolFactoryDefault()
	transportFactory := thrift.NewTBufferedTransportFactory(
		s.BardServerOptions.Info.BufferSize)
	s.Server = thrift.NewTSimpleServer4(processor, transport,
		transportFactory, protoFactory)
	err = s.Server.Serve()
	return
}

func (s *ThriftServer) Stop() {
	s.Server.Stop()
}