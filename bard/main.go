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
	"strings"
)

var logLevel string

var httpAddr string
var rpcAddr string

var barExe string

var storageType string

var serverInfo proto.ServerInfo
var rpcEndpoints string

var blockStorageOptions storage.BlockStorageOptions

func init() {
	flag.StringVar(&logLevel, "log-level", logx.ERROR, "logging level")

	flag.StringVar(&httpAddr, "bind-http", ":3000", "HTTP bind addr")
	flag.StringVar(&rpcAddr, "bind-rpc", ":3000", "RPC bind addr")


	flag.StringVar(&serverInfo.HTTPEndpoint, "http", "http://localhost:3000/v1",
		"HTTP endpoint")
	flag.StringVar(&rpcEndpoints, "rpc", "localhost:3001", "RPC endpoints")
	flag.Int64Var(&serverInfo.ChunkSize, "chunk", 1024 * 1024 * 2, "preferred chunk size")
	flag.IntVar(&serverInfo.PoolSize, "conns", 16,
		"preferred conns from one client")
	flag.IntVar(&serverInfo.BufferSize, "buffer", 1024 * 1024 * 8, "thrift buffer size")

	flag.StringVar(&barExe, "bar-exe", "",
		"path to windows bar executable")

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
	serverInfo.RPCEndpoints = strings.Split(rpcEndpoints, ",")
	srv := server.NewBardServer(&server.BardServerOptions{
		HttpBind: httpAddr,
		RPCBind: rpcAddr,
		ServerInfo: &serverInfo,
		Storage: pool,
		BarExe: barExe,
	})

	err = srv.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func newStorage(kind string) (res storage.Storage, err error) {
	res = storage.NewBlockStorage(&blockStorageOptions)
	return
}
