package proto_test
import (
	"testing"
	"net/rpc"
	"net"
	"github.com/akaspin/bar/fixtures"
	"fmt"
	"github.com/stretchr/testify/assert"
"net/http"
)


type TestMessage struct {
	Id string
	Data []byte
}

type ServiceFixture struct {}


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

type binaryTestServer struct {
	*baseTestServer
}

func (s *binaryTestServer) start() (err error) {

	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			continue
		}
		s.service.ServeConn(conn)
	}
	return
}

type httpTestServer struct {
	*baseTestServer
}

func (s *httpTestServer) start() (err error) {
	mux := http.NewServeMux()
	mux.Handle("/v1/rpc", s.service)
	srv := &http.Server{Handler:mux}
	err = srv.Serve(s.Listener)
	return
}

func Test_Proto_Binary1(t *testing.T)  {
	srv := &binaryTestServer{&baseTestServer{}}
	port, err := srv.listen()
	assert.NoError(t, err)
	go srv.start()
	defer srv.Stop()

	client, err := rpc.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	assert.NoError(t, err)
	defer client.Close()

	var res TestMessage
	req := TestMessage{
		"12345678910",
		[]byte("mama myla ramu"),
	}
	err = client.Call("ServiceFixture.Ping", &req, &res)
	assert.NoError(t, err)
	assert.Equal(t, req, res)
}


func Test_Proto_HTTP1(t *testing.T)  {
	srv := &httpTestServer{&baseTestServer{}}
	port, err := srv.listen()
	assert.NoError(t, err)
	go srv.start()
	defer srv.Stop()

	client, err := rpc.DialHTTPPath("tcp",
		fmt.Sprintf("127.0.0.1:%d", port), "/v1/rpc")
	assert.NoError(t, err)
	defer client.Close()

	var res TestMessage
	req := TestMessage{
		"12345678910",
		[]byte("mama myla ramu"),
	}
	err = client.Call("ServiceFixture.Ping", &req, &res)
	assert.NoError(t, err)
	assert.Equal(t, req, res)
}

func Benchmark_Proto_Binary1(b *testing.B)  {
	srv := &binaryTestServer{&baseTestServer{}}
	port, err := srv.listen()
	assert.NoError(b, err)
	go srv.start()
	defer srv.Stop()

	client, err := rpc.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	assert.NoError(b, err)
	defer client.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		var res TestMessage
		req := TestMessage{
			"12345678910",
			[]byte("mama myla ramu"),
		}
		err = client.Call("ServiceFixture.Ping", &req, &res)
		b.StopTimer()
		assert.NoError(b, err)
		assert.Equal(b, req, res)
	}
}

func Benchmark_Proto_HTTP1(b *testing.B)  {
	srv := &binaryTestServer{&baseTestServer{}}
	port, err := srv.listen()
	assert.NoError(b, err)
	go srv.start()
	defer srv.Stop()

	client, err := rpc.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	assert.NoError(b, err)
	defer client.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		var res TestMessage
		req := TestMessage{
			"12345678910",
			[]byte("mama myla ramu"),
		}
		err = client.Call("ServiceFixture.Ping", &req, &res)
		b.StopTimer()
		assert.NoError(b, err)
		assert.Equal(b, req, res)
	}
}
