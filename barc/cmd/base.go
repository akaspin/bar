package cmd
import (
	"flag"
	"io"
	"github.com/tamtam-im/logx"
	"github.com/akaspin/bar/proto"
	"github.com/tamtam-im/flags"
)

type Env struct {
	WD string
	Flags *flags.Flags
	Stdin io.Reader
	Stdout io.Writer
	Stderr io.Writer

}

type BaseOptions struct  {
	LogLevel string
	ChunkSize int64
	Endpoints string
	PoolSize int
	BufferSize int
}


// Subcommand environment
type Base struct {
	*BaseOptions
	*Env
}

func (b *Base) Init(fs *flag.FlagSet) {}

func NewBaseCmd(args []string, in io.Reader, stdout, stderr io.Writer) (res *Base) {
	res = &Base{
		BaseOptions: &BaseOptions{

		},
		Env: &Env{
			Flags: flags.New(flag.NewFlagSet(args[0], flag.ExitOnError)),
			Stdin: in,
			Stdout: stdout,
			Stderr: stderr,
		},
	}

	res.Flags.StringVar(&res.WD, "work-dir", "", "set work directory")
	res.Flags.StringVar(&res.LogLevel, "log-level", logx.INFO, "logging level")
	res.Flags.StringVar(&res.Endpoints, "endpoint", "localhost:3000",
		"bard endpoints separated by comma")
	res.Flags.Int64Var(&res.ChunkSize, "chunk", proto.CHUNK_SIZE, "preferred chunk size")
	res.Flags.IntVar(&res.PoolSize, "pool", 16, "preferred connection pool size")
	res.Flags.IntVar(&res.BufferSize, "buffer", 1024 * 1024 * 8, "thrift buffer size")
	return
}

func (b *Base) Description() string {
	return "bar desc"
}

func (b *Base) Help() {
	b.Flags.Help()
}

func (b *Base) Summary() {
	b.Flags.Summary()
}