package main
import (
	"flag"
	"github.com/tamtam-im/flags"
	"os"
	"github.com/akaspin/bar/bard/storage"
	"time"
	"fmt"
	"github.com/akaspin/bar/bard/server"
	"github.com/tamtam-im/logx"
	"github.com/akaspin/bar/proto"
)

var logLevel string

var httpAddr string
var rpcAddr string

var chunkSize int64
var clientConns int
var endpoint string
var httpEndpoint string
var barExe string

var storageType string
var storageWorkers int

var storageBlockRoot string
var storageBlockSplit int

func init() {
	flag.StringVar(&logLevel, "log-level", logx.ERROR, "logging level")

	flag.StringVar(&httpAddr, "bind-http", ":3001", "HTTP bind addr")
	flag.StringVar(&rpcAddr, "bind-rpc", "3000", "RPC bind addr")


	flag.Int64Var(&chunkSize, "chunk", 1024*1024*2, "preferred chunk size")
	flag.IntVar(&clientConns, "conns", 16,
		"preferred conns from one client")
	flag.StringVar(&endpoint, "endpoint", "http://localhost:3000/v1",
		"RPC endpoint")
	flag.StringVar(&httpEndpoint, "http-endpoint", "http://localhost:3000/v1",
		"HTTP endpoint")
	flag.StringVar(&barExe, "barc-exe", "",
		"path to windows barc executable")

	flag.StringVar(&storageType, "storage-type", "block", "storage type")
	flag.IntVar(&storageWorkers, "storage-workers", 128, "storage workers")

	// block storage options
	flag.StringVar(&storageBlockRoot, "storage-block-root", "data",
		"block storage root")
	flag.IntVar(&storageBlockSplit, "storage-block-split", 2,
		"block storage split factor")
}

func main() {
	flags.New(flag.CommandLine).Boot(os.Args)
	logx.SetLevel(logLevel)

	pool, err := storagePool(storageType)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	srv := server.NewBardServer(&server.BardServerOptions{
		httpAddr,
		rpcAddr,
		&proto.Info{httpEndpoint, endpoint, chunkSize, clientConns},
		pool,
		barExe,
	})

	err = srv.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func storagePool(kind string) (res *storage.StoragePool, err error) {
	res = storage.NewStoragePool(
		storage.NewBlockStorageFactory(storageBlockRoot, storageBlockSplit),
		storageWorkers, time.Minute * 5)

	return
}
