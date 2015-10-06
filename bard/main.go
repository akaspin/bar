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

var addr string
var chunkSize int64
var clientConns int
var endpoint string

var storageType string
var storageWorkers int

var storageBlockRoot string
var storageBlockSplit int

func init() {
	flag.StringVar(&logLevel, "logging-level", logx.ERROR, "logging level")

	flag.StringVar(&addr, "bind", ":3000", "bind addr")
	flag.Int64Var(&chunkSize, "chunk", 1024*1024*2, "preferred chunk size")
	flag.IntVar(&clientConns, "conns", 16,
		"preferred conns from one client")
	flag.StringVar(&endpoint, "endpoint", "http://localhost:3000/v1", "endpoint")

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
		addr,
		&proto.Info{endpoint, chunkSize, clientConns},
		pool,
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
