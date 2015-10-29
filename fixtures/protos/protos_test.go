package protos_test

import (
	"fmt"
	"github.com/akaspin/bar/fixtures"
	"net"
	"net/rpc"
)

type TestMessage struct {
	Id   string
	Data []byte
}

type ServiceFixture struct{}

func (s *ServiceFixture) Ping(req *TestMessage, res *TestMessage) (err error) {
	*res = *req
	return
}

func newRpc() (res *rpc.Server, err error) {
	res = rpc.NewServer()
	rpcService := &ServiceFixture{}
	err = res.Register(rpcService)
	return
}

type baseTestServer struct {
	net.Listener
	service *rpc.Server
}

func (s *baseTestServer) listen() (port int, err error) {
	if s.service, err = newRpc(); err != nil {
		return
	}
	if port, err = fixtures.GetOpenPort(); err != nil {
		return
	}
	s.Listener, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
	return
}

func (s *baseTestServer) Stop() {
	s.Listener.Close()
}
