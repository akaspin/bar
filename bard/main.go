package main
import (
	"flag"
	"github.com/tamtam-im/flags"
	"os"
	"github.com/akaspin/bar/bard/storage"
	"fmt"
	"github.com/akaspin/bar/bard/server"
	"github.com/tamtam-im/logx"
	"github.com/akaspin/bar/proto"
)

var logLevel string

var httpAddr string
var rpcAddr string

var chunkSize int64
var poolConns int

var rpcEndpoint string
var httpEndpoint string
var barExe string

var storageType string

var blockStorageOptions storage.BlockStorageOptions

func init() {
	flag.StringVar(&logLevel, "log-level", logx.ERROR, "logging level")

	flag.StringVar(&httpAddr, "bind-http", ":3000", "HTTP bind addr")
	flag.StringVar(&rpcAddr, "bind-rpc", ":3000", "RPC bind addr")


	flag.Int64Var(&chunkSize, "chunk", 1024*1024*2, "preferred chunk size")
	flag.IntVar(&poolConns, "conns", 16,
		"preferred conns from one client")
	flag.StringVar(&httpEndpoint, "http", "http://localhost:3000/v1",
		"HTTP endpoint")
	flag.StringVar(&rpcEndpoint, "rpc", "http://localhost:3000/v1",
		"RPC endpoint")
	flag.StringVar(&barExe, "barc-exe", "",
		"path to windows barc executable")

	flag.StringVar(&storageType, "storage-type", "block", "storage type")

	// block storage options
	flag.StringVar(&blockStorageOptions.Root, "storage-block-root", "data",
		"block storage root")
	flag.IntVar(&blockStorageOptions.Split, "storage-block-split", 2,
		"block storage split factor")
	flag.IntVar(&blockStorageOptions.MaxFiles, "storage-block-max-files", 32,
		"block storage max open files (ulimit -n)")
	flag.IntVar(&blockStorageOptions.PoolSize, "storage-block-pool", 32,
		"block storage pool size")
}

func main() {
	flags.New(flag.CommandLine).Boot(os.Args)
	logx.SetLevel(logLevel)

	pool, err := newStorage(storageType)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	srv, err := server.NewBardServer(&server.BardServerOptions{
		HttpBind: httpAddr,
		RPCBind: rpcAddr,
		Info: &proto.Info{
			HTTPEndpoint: httpEndpoint,
			RPCEndpoints: []string{rpcEndpoint},
			ChunkSize: chunkSize,
			PoolSize: poolConns,
			BufferSize: 1024 * 1024 * 8,
		},
		Storage: pool,
		BarExe: barExe,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	err = srv.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func newStorage(kind string) (res storage.Storage, err error) {
	res = storage.NewBlockStorage(&blockStorageOptions)
	return
}
