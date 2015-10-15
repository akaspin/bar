package server

import (
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/akaspin/bar/proto/wire"
	"github.com/tamtam-im/logx"
)

type ThriftServer struct  {
	*BardServerOptions
	Server thrift.TServer
}

func NewThriftServer(options *BardServerOptions) *ThriftServer  {
	return &ThriftServer{BardServerOptions: options}
}

func (s *ThriftServer) Start() (err error) {
	handler := NewBardTHandler(s.BardServerOptions)
	processor := wire.NewBarProcessor(handler)

	var transport thrift.TServerTransport

	if transport, err = thrift.NewTServerSocket(s.RPCBind); err != nil {
		return
	}
	protoFactory := thrift.NewTBinaryProtocolFactoryDefault()
	transportFactory := thrift.NewTBufferedTransportFactory(
		s.BardServerOptions.ServerInfo.BufferSize)
	s.Server = thrift.NewTSimpleServer4(processor, transport,
		transportFactory, protoFactory)

	logx.Debugf("thrift listening at %s", s.RPCBind)
	err = s.Server.Serve()
	return
}

func (s *ThriftServer) Stop() {
	s.Server.Stop()
	logx.Debugf("thrift stopped at %s", s.RPCBind)
}