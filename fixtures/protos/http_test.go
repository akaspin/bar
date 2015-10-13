package protos_test
import (
"net/http"
"testing"
	"github.com/stretchr/testify/assert"
	"net/rpc"
	"fmt"
	"time"
	"strings"
)

type httpTestServer struct {
	*baseTestServer
}

func (s *httpTestServer) start() (err error) {
	mux := http.NewServeMux()
	mux.Handle("/v1/rpc", s.service)
	srv := &http.Server{Handler:mux}
	go srv.Serve(s.Listener)
//	time.Sleep(time.Second)
	return
}

func Test_Proto_HTTP(t *testing.T)  {
	srv := &httpTestServer{&baseTestServer{}}
	port, err := srv.listen()
	assert.NoError(t, err)
	srv.start()
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

func Test_Proto_HTTP_Large(t *testing.T)  {
	srv := &httpTestServer{&baseTestServer{}}
	port, err := srv.listen()
	assert.NoError(t, err)
	srv.start()
	defer srv.Stop()

	client, err := rpc.DialHTTPPath("tcp",
		fmt.Sprintf("127.0.0.1:%d", port), "/v1/rpc")
	assert.NoError(t, err)
	defer client.Close()

	var res TestMessage
	req := TestMessage{
		"12345678910",
		[]byte(strings.Repeat("0", 1024 * 1024)),
	}
	err = client.Call("ServiceFixture.Ping", &req, &res)
	assert.NoError(t, err)
	assert.Equal(t, req, res)
}

func Test_Proto_HTTP_Timeout(t *testing.T)  {
	srv := &httpTestServer{&baseTestServer{}}
	port, err := srv.listen()
	assert.NoError(t, err)
	srv.start()
	defer srv.Stop()

	client, err := rpc.DialHTTPPath("tcp",
		fmt.Sprintf("127.0.0.1:%d", port), "/v1/rpc")
	assert.NoError(t, err)
	defer client.Close()

	var res, res1 TestMessage
	req := TestMessage{
		"12345678910",
		[]byte("mama myla ramu"),
	}
	err = client.Call("ServiceFixture.Ping", &req, &res)
	assert.NoError(t, err)
	assert.Equal(t, req, res)

	time.Sleep(time.Second)
	err = client.Call("ServiceFixture.Ping", &req, &res1)
	assert.NoError(t, err)
	assert.Equal(t, req, res1)
}

func Test_Proto_HTTP_Large_Timeout(t *testing.T)  {
	srv := &httpTestServer{&baseTestServer{}}
	port, err := srv.listen()
	assert.NoError(t, err)
	srv.start()
	defer srv.Stop()

	client, err := rpc.DialHTTPPath("tcp",
		fmt.Sprintf("127.0.0.1:%d", port), "/v1/rpc")
	assert.NoError(t, err)
	defer client.Close()

	var res, res1 TestMessage
	req := TestMessage{
		"12345678910",
		[]byte(strings.Repeat("0", 1024 * 1024)),
	}
	err = client.Call("ServiceFixture.Ping", &req, &res)
	assert.NoError(t, err)
	assert.Equal(t, req, res)

	time.Sleep(time.Second)
	err = client.Call("ServiceFixture.Ping", &req, &res1)
	assert.NoError(t, err)
	assert.Equal(t, req, res1)
}

func Benchmark_Proto_HTTP(b *testing.B)  {
	b.StopTimer()
	srv := &httpTestServer{&baseTestServer{}}
	port, err := srv.listen()
	assert.NoError(b, err)
	srv.start()
	defer srv.Stop()

	client, err := rpc.DialHTTPPath("tcp",
		fmt.Sprintf("127.0.0.1:%d", port), "/v1/rpc")
	assert.NoError(b, err)
	defer client.Close()

	b.Log(port, b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
			b.StartTimer()
			var res TestMessage
			req := TestMessage{
				fmt.Sprintf("%d", i),
				[]byte("mama myla ramu"),
			}
			err = client.Call("ServiceFixture.Ping", &req, &res)
			b.StopTimer()
			assert.NoError(b, err)
			assert.Equal(b, req, res)
			b.SetBytes(int64(len(req.Data) * 2))
	}
}

func Benchmark_Proto_HTTP_Large(b *testing.B)  {
	srv := &httpTestServer{&baseTestServer{}}
	port, err := srv.listen()
	assert.NoError(b, err)
	srv.start()
	defer srv.Stop()

	client, err := rpc.DialHTTPPath("tcp",
		fmt.Sprintf("127.0.0.1:%d", port), "/v1/rpc")
	assert.NoError(b, err)
	defer client.Close()

	b.Log(port, b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
			b.StartTimer()
			var res TestMessage
			req := TestMessage{
				fmt.Sprintf("%d", i),
				[]byte(strings.Repeat("0", 1024 * 1024)),
			}
			err = client.Call("ServiceFixture.Ping", &req, &res)
			b.StopTimer()
			assert.NoError(b, err)
			assert.Equal(b, req, res)
			b.SetBytes(int64(len(req.Data) * 2))
	}
}
