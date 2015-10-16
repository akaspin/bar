package cmd
import (
	"strings"
	"flag"
	"io"
	"github.com/tamtam-im/logx"
	"github.com/akaspin/bar/proto"
"github.com/tamtam-im/flags"
)

type BaseOptions struct  {
	WD string
	LogLevel string
	ChunkSize int64
	PoolSize int
	BufferSize int

	endpoints string
}

func (o *BaseOptions) Endpoints() []string {
	return strings.Split(o.endpoints, ",")
}

// Subcommand environment
type Base struct {
	WD string
	*BaseOptions

	FS *flag.FlagSet
	Stdin io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func (b *Base) Init(fs *flag.FlagSet) {}

func NewBaseCmd(args []string, in io.Reader, stdout, stderr io.Writer) (res *Base) {
	res = &Base{
		BaseOptions: &BaseOptions{},
		FS: flag.NewFlagSet(args[0], flag.ExitOnError),
		Stdin: in,
		Stdout: stdout,
		Stderr: stderr,
	}

	res.FS.StringVar(&res.WD, "work-dir", "", "set work directory")
	res.FS.StringVar(&res.LogLevel, "log-level", logx.INFO, "logging level")
	res.FS.StringVar(&res.endpoints, "endpoint", "localhost:3000",
		"bard endpoints separated by comma")
	res.FS.Int64Var(&res.ChunkSize, "chunk", proto.CHUNK_SIZE, "preferred chunk size")
	res.FS.IntVar(&res.PoolSize, "pool", 16, "preferred connection pool size")
	res.FS.IntVar(&res.BufferSize, "buffer", 1024 * 1024 * 8, "thrift buffer size")
	return
}

func (b *Base) Description() string {
	return "BAR"
}

func (b *Base) Help() (res string) {
	res = flags.DEFAULT_HELP
	return
}

func (b *Base) Summary() (res string) {
	res = flags.DEFAULT_SUMMARY
	return
}