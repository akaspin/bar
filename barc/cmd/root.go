package cmd

import (
	"flag"
	"fmt"
	"github.com/tamtam-im/flags"
	"github.com/tamtam-im/logx"
	"io"
)

var logLevel string

type SubCommand interface {

	// Do subcommand
	Do() (err error)
}

type SubCommandFactory func(*BaseSubCommand) SubCommand

// Subcommand environment
type BaseSubCommand struct {
	WD string
	FS *flag.FlagSet
	Stdin io.Reader
	Stdout io.Writer
	StdErr io.Writer
}

func route(s string, base *BaseSubCommand) (res SubCommand, err error) {
	factory, ok := (map[string]SubCommandFactory{
//		"git-init":       NewGitInitCmd,
		"git-clean":      NewGitCleanCommand,
		"git-smudge":     NewGitSmudgeCmd,
//		"git-pre-commit": NewGitPreCommitCmd,
//		"up":             NewUpCmd,
//		"down":           NewDownCmd,
//		"ls":             NewLsCmd,
		"git-diff":       NewGitDiffCmd,
//		"spec-out":       NewSpecOutCmd,
	})[s]
	if !ok {
		err = fmt.Errorf("command %s not found", s)
	}

	return factory(base), nil
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
	subFS := flag.NewFlagSet(f.FlagSet.Args()[0], flag.ExitOnError)

	sub, err := route(f.FlagSet.Args()[0], &BaseSubCommand{
		wd, subFS, in, out, errOut,
	})
	if err != nil {
		return
	}

	flags.New(subFS).Boot(f.FlagSet.Args())

	logx.Debugf("invoking %s in %s with %s",
		f.FlagSet.Args()[0], wd, f.FlagSet.Args())
	err = sub.Do()
	return
}
