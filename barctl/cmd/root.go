package cmd
import (
	"github.com/tamtam-im/flags"
	"io"
	"flag"
	"fmt"
)

type SubCommand interface  {

	// Init flagset
	Bind(fs *flag.FlagSet, in io.Reader, out, errOut io.Writer) (err error)

	// Handle
	Do() (err error)
}

func route(s string) (res SubCommand, err error) {
	res, ok := (map[string]SubCommand{
		"git-clean": &GitCleanCommand{},
		"git-cat": &GitCatCommand{},
		"git-pre-commit": &GitPreCommitCmd{},
		"upload": &UploadCommand{},
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
	if err = sub.Bind(subFS, in, out, errOut); err != nil {
		return
	}

	flags.New(subFS).NoEnv().Boot(f.FlagSet.Args())
	return sub.Do()
}
