package cmd
import (
	"github.com/tamtam-im/flags"
	"io"
	"flag"
	"fmt"
	"github.com/tamtam-im/logx"
)

var logLevel string

type SubCommand interface  {

	// Bind subcomand to environment
	Bind(wd string, fs *flag.FlagSet, in io.Reader, out io.Writer) (err error)

	// Do subcommand
	Do() (err error)
}

func route(s string) (res SubCommand, err error) {
	res, ok := (map[string]SubCommand{
		"git-init": &GitInitCmd{},
		"git-clean": &GitCleanCommand{},
		"git-smudge": &GitSmudgeCmd{},
		"git-cat": &GitCatCommand{},
		"git-pre-commit": &GitPreCommitCmd{},
		"up": &UpCmd{},
		"down": &DownCmd{},
		"ls": &LsCmd{},
	})[s]
	if !ok {
		err = fmt.Errorf("command %s not found", s)
	}

	return
}

func Root(wd string, args []string, in io.Reader, out, errOut io.Writer) (err error) {
	flag.StringVar(&logLevel, "log-level", logx.DEBUG, "logging level")

	f := flags.New(flag.CommandLine)
	f.Boot(args)

	logx.SetLevel(logLevel)
	logx.SetOutput(errOut)

	// route subcommand
	if len(f.FlagSet.Args()) == 0 {
		f.Usage()
	}
	sub, err := route(f.FlagSet.Args()[0])
	if err != nil {
		return
	}

	subFS := flag.NewFlagSet(f.FlagSet.Args()[0], flag.ExitOnError)
	if err = sub.Bind(wd, subFS, in, out); err != nil {
		return
	}

	flags.New(subFS).Boot(f.FlagSet.Args())

	logx.Debugf("invoking %s in %s with %s",
		f.FlagSet.Args()[0], wd, f.FlagSet.Args())
	err = sub.Do()
	return
}
