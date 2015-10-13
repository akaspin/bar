package pb3_test
import (
	"net"
	"golang.org/x/net/context"
	"github.com/akaspin/bar/fixtures/protos/pb3"
	"github.com/tamtam-im/logx"
	"google.golang.org/grpc"
	"testing"
	"fmt"
	"github.com/akaspin/bar/fixtures"
	"github.com/stretchr/testify/assert"
	"strings"
	"time"
)


type server struct {}

func (s *server) Test(ctx context.Context, req *pb3.TestRep) (res *pb3.TestRep, err error) {
	res = req
	return
}

type GRPCServer struct {
	l net.Listener
	srv *grpc.Server
	Port int
	err error
}

func (s *GRPCServer) Start() {
	if s.Port, s.err = fixtures.GetOpenPort(); s.err != nil {
		return
	}

	if s.l, s.err = net.Listen("tcp", fmt.Sprintf(":%d", s.Port)); s.err != nil {
		return
	}

	s.srv = grpc.NewServer()
	pb3.RegisterBarServer(s.srv, &server{})
	go s.srv.Serve(s.l)
	time.Sleep(time.Second)
}


func Test_Proto_GRPC(t *testing.T) {
	logx.SetLevel(logx.DEBUG)
	srv1 := new(GRPCServer)
	srv1.Start()
	defer srv1.srv.Stop()

	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", srv1.Port), grpc.WithInsecure())
	assert.NoError(t, err)
	defer conn.Close()

	c := pb3.NewBarClient(conn)

	req := pb3.TestRep{
		"test",
		[]byte(strings.Repeat("m", 10)),
	}
	res, err := c.Test(context.Background(), &req)
	assert.NoError(t, err)
	assert.EqualValues(t, req, *res)
}

func Test_Proto_GRPC_Timeout(t *testing.T) {
	logx.SetLevel(logx.DEBUG)
	srv1 := new(GRPCServer)
	srv1.Start()
	defer srv1.srv.Stop()

	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", srv1.Port), grpc.WithInsecure())
	assert.NoError(t, err)
	defer conn.Close()

	c := pb3.NewBarClient(conn)

	req := pb3.TestRep{
		"test",
		[]byte(strings.Repeat("m", 10)),
	}
	res, err := c.Test(context.Background(), &req)
	assert.NoError(t, err)
	assert.EqualValues(t, req, *res)

	res1, err := c.Test(context.Background(), &req)
	assert.NoError(t, err)
	assert.EqualValues(t, req, *res1)
}

func Test_Proto_GRPC_Large(t *testing.T) {
	logx.SetLevel(logx.DEBUG)
	srv1 := new(GRPCServer)
	srv1.Start()
	defer srv1.srv.Stop()

	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", srv1.Port), grpc.WithInsecure())
	assert.NoError(t, err)
	defer conn.Close()

	c := pb3.NewBarClient(conn)

	req := pb3.TestRep{
		"test",
		[]byte(strings.Repeat("m", 1024 * 1024)),
	}
	res, err := c.Test(context.Background(), &req)
	assert.NoError(t, err)
	assert.EqualValues(t, req, *res)
}

func Benchmark_Proto_GRPC(b *testing.B) {
	logx.SetLevel(logx.DEBUG)
	srv1 := new(GRPCServer)
	srv1.Start()
	defer srv1.srv.Stop()

	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", srv1.Port), grpc.WithInsecure())
	assert.NoError(b, err)
	defer conn.Close()

	c := pb3.NewBarClient(conn)

	b.Log(srv1.Port, b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
			b.StartTimer()
			req := pb3.TestRep{
				fmt.Sprintf("%d", i),
				[]byte("mama myla ramu"),
			}
			res, err := c.Test(context.Background(), &req)
			b.StopTimer()
			assert.NoError(b, err)
			assert.EqualValues(b, req, *res)
			b.SetBytes(int64(len(req.Body) * 2))
	}
}

func Benchmark_Proto_GRPC_Large(b *testing.B) {
	logx.SetLevel(logx.DEBUG)
	srv1 := new(GRPCServer)
	srv1.Start()
	defer srv1.srv.Stop()

	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", srv1.Port), grpc.WithInsecure())
	assert.NoError(b, err)
	defer conn.Close()

	c := pb3.NewBarClient(conn)

	b.Log(srv1.Port, b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
			b.StartTimer()
			req := pb3.TestRep{
				fmt.Sprintf("%d", i),
				[]byte(strings.Repeat("m", 1024 * 1024)),
			}
			res, err := c.Test(context.Background(), &req)
			b.StopTimer()
			assert.NoError(b, err)
			assert.EqualValues(b, req, *res)
			b.SetBytes(int64(len(req.Body) * 2))
	}
}
