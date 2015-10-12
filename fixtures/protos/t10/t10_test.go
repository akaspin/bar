package t10_test
import (
	"github.com/akaspin/bar/fixtures/protos/t10/srv"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/akaspin/bar/fixtures"
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
"strings"
	"time"
)

type service struct {}

func (s *service) Test(par *srv.TestRep) (r *srv.TestRep, err error) {
	r = par
	return
}

type TServer struct {
	Port int
	Server thrift.TServer
}

func (s *TServer) Start() (err error) {
	if s.Port, err = fixtures.GetOpenPort(); err != nil {
		return
	}

	protoFactory := thrift.NewTBinaryProtocolFactoryDefault()
	transportFactory := thrift.NewTTransportFactory()
	processor := srv.NewTSrvProcessor(&service{})

	var transport thrift.TServerTransport
	if transport, err = thrift.NewTServerSocket(fmt.Sprintf(":%d", s.Port)); err != nil {
		return
	}
	s.Server = thrift.NewTSimpleServer4(processor, transport,
		transportFactory, protoFactory)
	go s.Server.Serve()
	time.Sleep(time.Second)
	return
}


func Test_Proto_Thrift(t *testing.T) {
	server := &TServer{}
	err := server.Start()
	assert.NoError(t, err)
	defer server.Server.Stop()

	var transport thrift.TTransport
	transport, err = thrift.NewTSocket(fmt.Sprintf("127.0.0.1:%d", server.Port))
	assert.NoError(t, err)

	transportFactory := thrift.NewTTransportFactory()
	transport = transportFactory.GetTransport(transport)

	err = transport.Open()
	assert.NoError(t, err)
	defer transport.Close()

	protoFactory := thrift.NewTBinaryProtocolFactoryDefault()

	client := srv.NewTSrvClientFactory(transport, protoFactory)

	req := srv.TestRep{
		ID: "test",
		Data: []byte("mama myla ramy"),
	}
	res, err := client.Test(&req)
	assert.NoError(t, err)
	assert.EqualValues(t, req, *res)
}

func Benchmark_Proto_Thrift(t *testing.B) {
	server := &TServer{}
	err := server.Start()
	assert.NoError(t, err)
	defer server.Server.Stop()

	var transport thrift.TTransport
	transport, err = thrift.NewTSocket(fmt.Sprintf("127.0.0.1:%d", server.Port))
	assert.NoError(t, err)

	transportFactory := thrift.NewTTransportFactory()
	transport = transportFactory.GetTransport(transport)

	err = transport.Open()
	assert.NoError(t, err)
	defer transport.Close()

	protoFactory := thrift.NewTBinaryProtocolFactoryDefault()

	client := srv.NewTSrvClientFactory(transport, protoFactory)

	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		t.StartTimer()
		req := srv.TestRep{
			ID: "test",
			Data: []byte("mama myla ramy"),
		}
		res, err := client.Test(&req)
		t.StopTimer()
		assert.NoError(t, err)
		assert.EqualValues(t, req, *res)
	}
}

func Benchmark_Proto_Thrift_Large(t *testing.B) {
	server := &TServer{}
	err := server.Start()
	assert.NoError(t, err)
	defer server.Server.Stop()

	var transport thrift.TTransport
	transport, err = thrift.NewTSocket(fmt.Sprintf("127.0.0.1:%d", server.Port))
	assert.NoError(t, err)

	transportFactory := thrift.NewTTransportFactory()
	transport = transportFactory.GetTransport(transport)

	err = transport.Open()
	assert.NoError(t, err)
	defer transport.Close()

	protoFactory := thrift.NewTBinaryProtocolFactoryDefault()

	client := srv.NewTSrvClientFactory(transport, protoFactory)

	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		t.StartTimer()
		req := srv.TestRep{
			ID: "test",
			Data: []byte(strings.Repeat("0", 1024 * 1024)),
		}
		res, err := client.Test(&req)
		t.StopTimer()
		assert.NoError(t, err)
		assert.EqualValues(t, req, *res)
	}
}