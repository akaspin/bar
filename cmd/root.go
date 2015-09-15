package cmd
import (
	"github.com/tamtam-im/flags"
	"os"
	"github.com/tamtam-im/logx"
	"io"
)

type subcommand func(args []string, in io.Reader, out, errOut io.Writer) error

func init() {
	// Common flags

}

func requireStdin(msg string) {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		logx.Fatal("Cannot read from STDIN. %s", msg)
		os.Exit(1)
	}
}

func Root(args []string, in io.Reader, out, errOut io.Writer) error {
	f := flags.New()
	f.Boot()

	// route subcommand
	if len(f.FlagSet.Args()) == 0 {
		f.FlagSet.Usage()
	}

	return (map[string]subcommand{
		"clean": CleanCmd,
	})[f.FlagSet.Args()[0]](f.FlagSet.Args(), in, out, errOut)
}
