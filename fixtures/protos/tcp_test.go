package protos_test
import (
	"testing"
	"github.com/stretchr/testify/assert"
	"net/rpc"
	"fmt"
	"time"
	"strings"
)

type binaryTestServer struct {
	*baseTestServer
}

func (s *binaryTestServer) start() (err error) {

	go func() {
		for {
			conn, err := s.Listener.Accept()
			if err != nil {
				continue
			}
			go s.service.ServeConn(conn)
		}
	}()
	time.Sleep(time.Second)
	return
}

func Test_Proto_Binary(t *testing.T)  {
	t.Skip()
	srv := &binaryTestServer{&baseTestServer{}}
	port, err := srv.listen()
	assert.NoError(t, err)
	srv.start()
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

func Test_Proto_Binary_Large(t *testing.T)  {
	t.Skip()
	srv := &binaryTestServer{&baseTestServer{}}
	port, err := srv.listen()
	assert.NoError(t, err)
	srv.start()
	defer srv.Stop()

	client, err := rpc.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
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

func Test_Proto_Binary_Large_Timeout(t *testing.T)  {
	t.Skip()
	srv := &binaryTestServer{&baseTestServer{}}
	port, err := srv.listen()
	assert.NoError(t, err)
	srv.start()
	defer srv.Stop()

	client, err := rpc.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
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

func Test_Proto_Binary_Timeout(t *testing.T)  {
	t.Skip()
	srv := &binaryTestServer{&baseTestServer{}}
	port, err := srv.listen()
	assert.NoError(t, err)
	srv.start()
	defer srv.Stop()

	client, err := rpc.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
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

func Benchmark_Proto_Binary(b *testing.B)  {
	b.Skip()
	srv := &binaryTestServer{&baseTestServer{}}
	port, err := srv.listen()
	assert.NoError(b, err)
	srv.start()
	defer srv.Stop()

	client, err := rpc.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	assert.NoError(b, err)
	defer client.Close()

	b.Log(port)

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
			b.SetBytes(int64(len(req.Data) * 2))
	}
}

func Benchmark_Proto_Binary_Large(b *testing.B)  {
	b.Skip()
	srv := &binaryTestServer{&baseTestServer{}}
	port, err := srv.listen()
	assert.NoError(b, err)
	srv.start()
	defer srv.Stop()

	client, err := rpc.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	assert.NoError(b, err)
	defer client.Close()

	b.Log(port)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
			b.StartTimer()
			var res TestMessage
			req := TestMessage{
				"12345678910",
				[]byte(strings.Repeat("0", 1024 * 1024)),
			}
			err = client.Call("ServiceFixture.Ping", &req, &res)
			b.StopTimer()
			assert.NoError(b, err)
			assert.Equal(b, req, res)
			b.SetBytes(int64(len(req.Data) * 2))
	}
}


