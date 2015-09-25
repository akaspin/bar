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
)

var logLevel string

var addr string
var storageType string
var storageWorkers int

var storageBlockRoot string
var storageBlockSplit int

func init() {
	flag.StringVar(&logLevel, "logging-level", logx.ERROR, "logging level")

	flag.StringVar(&addr, "bind", ":3000", "bind addr")
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
	err = server.Serve(addr, pool)
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
