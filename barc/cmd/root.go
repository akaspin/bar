package cmd

import (
	"flag"
	"github.com/tamtam-im/flags"
	"github.com/tamtam-im/logx"
	"path/filepath"
)

type SubCommand interface {

	Init(fs *flag.FlagSet)

	// Do subcommand
	Do(args []string) (err error)

	Description() string

	Help() string
	Summary() string
}

type SubCommandFactory func(*Base) SubCommand

type RootCmd struct {
	*Base
}

func NewRootCmd(b *Base) (res *RootCmd, err error) {
	res = &RootCmd{b}
	return
}

func (r *RootCmd) Do(wd string, args []string) (err error) {
	f := flags.New(r.Base.FS).
		SetPrefix("BAR").
		SetSummary(r.Summary()).
		SetHelp(r.Help())
	f.Boot(args)

	if r.Base.WD != "" {
		r.Base.WD = filepath.Clean(filepath.Join(wd, r.Base.WD))
	}

	logx.SetLevel(r.Base.LogLevel)
	logx.SetOutput(r.Stderr)

	if len(r.FS.Args()) == 0 {
		r.FS.Usage()
	}
	subFactory, ok := r.getRoute()[r.FS.Args()[0]]
	if !ok {
		r.FS.Usage()
	}

	sub := subFactory(r.Base)
	subFS := flag.NewFlagSet(r.FS.Args()[0], flag.ExitOnError)
	sub.Init(subFS)

	flags.New(subFS).
		SetPrefix("BAR").
		SetHelp(sub.Help()).
		SetSummary(sub.Summary()).
		Boot(f.FlagSet.Args())

	logx.Debugf("invoking %s in %s with %s",
		f.FlagSet.Args()[0], wd, f.FlagSet.Args())
	err = sub.Do(f.FlagSet.Args()[1:])
	return
}

func (r *RootCmd) Help() string {
	return flags.DEFAULT_HELP
}

func (r *RootCmd) Usage() string {
	return flags.DEFAULT_SUMMARY
}

func (r *RootCmd) getRoute() map[string]SubCommandFactory {
	return map[string]SubCommandFactory{
		"git-init":       NewGitInitCmd,
		"git-clean":      NewGitCleanCommand,
		"git-smudge":     NewGitSmudgeCmd,
		"git-pre-commit": NewGitPreCommitCmd,
		"up":             NewUpCmd,
		"down":           NewDownCmd,
		"ls":             NewLsCmd,
		"spec-export":    NewSpecExportCmd,
		"spec-import":    NewSpecImportCmd,
	}
}
