package cmd
import (
	"github.com/tamtam-im/flags"
	"io"
	"flag"
	"fmt"
)

type SubCommand interface  {

	// Init flagset
	FS(fs *flag.FlagSet)

	//

	// Handle
	Do(in io.Reader, out, errOut io.Writer) (err error)
}

func addEndpointToFS(fs *flag.FlagSet, v *string)  {
	fs.StringVar(v, "endpoint", "http://127.0.0.1:3000/v1",
		"bard endpoint")
}


func route(s string) (res SubCommand, err error) {
	res, ok := (map[string]SubCommand{
		"git-clean": &GitCleanCommand{},
		"git-cat": &GitCatCommand{},
	})[s]
	if !ok {
		err = fmt.Errorf("%s not found")
	}

	return
}

func Root(args []string, in io.Reader, out, errOut io.Writer) (err error) {
	f := flags.New(flag.CommandLine).NoEnv()
	f.Boot(args)

	// route subcommand
	if len(f.FlagSet.Args()) == 0 {
		f.Usage()
	}
	sub, err := route(f.FlagSet.Args()[0])
	if err != nil {
		return
	}

	subFS := flag.NewFlagSet(f.FlagSet.Args()[0], flag.ExitOnError)
	sub.FS(subFS)

	flags.New(subFS).NoEnv().Boot(f.FlagSet.Args())
	return sub.Do(in, out, errOut)
}
