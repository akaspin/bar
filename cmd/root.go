package cmd
import (
	"github.com/tamtam-im/flags"
	"io"
	"flag"
)

type SubCommand interface  {

	// Init flagset
	FS(fs *flag.FlagSet)

	// Handle
	Do(in io.Reader, out, errOut io.Writer) (err error)
}

func addEndpointToFS(fs *flag.FlagSet, v *string)  {
	fs.StringVar(v, "endpoint", "http://127.0.0.1:3000/v1",
		"bard endpoint")
}


type subcommand func(args []string, in io.Reader, out, errOut io.Writer) error

func route(s string) SubCommand {
	return (map[string]SubCommand{
		"clean": &CleanSubCommand{},
	})[s]
}

func Root(args []string, in io.Reader, out, errOut io.Writer) (err error) {
	f := flags.New(flag.CommandLine).NoEnv()
	f.Boot(args)

	// route subcommand
	if len(f.FlagSet.Args()) == 0 {
		f.Usage()
	}

	sub := route(f.FlagSet.Args()[0])
	subFS := flag.NewFlagSet(f.FlagSet.Args()[0], flag.ExitOnError)
	sub.FS(subFS)

	flags.New(subFS).NoEnv().Boot(f.FlagSet.Args())
	return sub.Do(in, out, errOut)
}
