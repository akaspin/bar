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

const thrift_buffer = 1024 * 1024 * 8

func (s *service) Test(par *srv.TestRep) (r *srv.TestRep, err error) {
	r = par
	return
}

type TServer struct {
	Port int
	Server thrift.TServer
}

func (s *TServer) Start(transportFactory thrift.TTransportFactory) (err error) {
	if s.Port, err = fixtures.GetOpenPort(); err != nil {
		return
	}

	protoFactory := thrift.NewTBinaryProtocolFactoryDefault()
	processor := srv.NewTSrvProcessor(&service{})

	var transport thrift.TServerTransport
	if transport, err = thrift.NewTServerSocket(fmt.Sprintf(":%d", s.Port)); err != nil {
		return
	}
	s.Server = thrift.NewTSimpleServer4(processor, transport,
		transportFactory, protoFactory)
	go s.Server.Serve()
	time.Sleep(time.Millisecond * 100)
	return
}


func Test_Proto_Thrift(t *testing.T) {
	t.Skip()
	server := &TServer{}
	err := server.Start(thrift.NewTTransportFactory())
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

func Benchmark_Proto_Thrift_Buffered(b *testing.B) {
	b.Skip()
	server := &TServer{}
	err := server.Start(thrift.NewTBufferedTransportFactory(thrift_buffer))
	assert.NoError(b, err)
	defer server.Server.Stop()

	var transport thrift.TTransport
	transport, err = thrift.NewTSocket(fmt.Sprintf("127.0.0.1:%d", server.Port))
	assert.NoError(b, err)

	transportFactory := thrift.NewTBufferedTransportFactory(thrift_buffer)
	transport = transportFactory.GetTransport(transport)

	err = transport.Open()
	assert.NoError(b, err)
	defer transport.Close()

	protoFactory := thrift.NewTBinaryProtocolFactoryDefault()
	client := srv.NewTSrvClientFactory(transport, protoFactory)

	b.Log(server.Port, b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
			b.StartTimer()
			req := srv.TestRep{
				ID: fmt.Sprintf("%d", i),
				Data: []byte("mama myla ramy"),
			}
			res, err := client.Test(&req)
			b.StopTimer()
			assert.NoError(b, err)
			assert.EqualValues(b, req, *res)
			b.SetBytes(int64(len(req.Data) * 2))
	}
}

func Benchmark_Proto_Thrift_Buffered_Large(b *testing.B) {
	b.Skip()
	server := &TServer{}
	err := server.Start(thrift.NewTBufferedTransportFactory(thrift_buffer))
	assert.NoError(b, err)
	defer server.Server.Stop()

	var transport thrift.TTransport
	transport, err = thrift.NewTSocket(fmt.Sprintf("127.0.0.1:%d", server.Port))
	assert.NoError(b, err)

	transportFactory := thrift.NewTBufferedTransportFactory(thrift_buffer)
	transport = transportFactory.GetTransport(transport)

	err = transport.Open()
	assert.NoError(b, err)
	defer transport.Close()

	protoFactory := thrift.NewTBinaryProtocolFactoryDefault()
	client := srv.NewTSrvClientFactory(transport, protoFactory)

	b.Log(server.Port, b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
			b.StartTimer()
			req := srv.TestRep{
				ID: fmt.Sprintf("%d", i),
				Data: []byte(strings.Repeat("0", 1024 * 1024)),
			}
			res, err := client.Test(&req)
			b.StopTimer()
			assert.NoError(b, err)
			assert.EqualValues(b, req, *res)
			b.SetBytes(int64(len(req.Data) * 2))
	}
}
