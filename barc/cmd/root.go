package cmd

import (
	"flag"
	"github.com/tamtam-im/flags"
	"github.com/tamtam-im/logx"
	"path/filepath"
	"fmt"
	"sort"
)

type SubCommand interface {

	Init(fs *flag.FlagSet)

	// Do subcommand
	Do(args []string) (err error)

	Description() string
	Help()
	Summary()
}

type SubCommandFactory func(*Base) SubCommand

type RootCmd struct {
	*Base
	baseSummary func()
}

func NewRootCmd(b *Base) (res *RootCmd, err error) {
	res = &RootCmd{Base: b}
	return
}

func (r *RootCmd) Do(wd string, args []string) (err error) {
	r.Base.Flags.Help = r.Help
	r.Base.Flags.Summary = r.Summary

	if res, err := r.Base.Flags.Boot(args); err != nil || res() {
		return err
	}

	if r.Base.WD != "" {
		r.Base.WD = filepath.Clean(filepath.Join(wd, r.Base.WD))
	}

	logx.SetLevel(r.Base.LogLevel)
	logx.SetOutput(r.Stderr)

	if len(r.Flags.Args()) == 0 {
		err = fmt.Errorf("no subcommand")
		r.Flags.Usage()
	}
	subFactory, ok := r.getRoute()[r.Flags.Args()[0]]
	if !ok {
		err = fmt.Errorf("invalid subcommand")
		r.Flags.Usage()
	}

	sub := subFactory(r.Base)
	subFS := flag.NewFlagSet(r.Flags.Args()[0], flag.ExitOnError)
	sub.Init(subFS)

	subFlags := flags.New(subFS).SetPrefix("BAR")
	subFlags.Help = sub.Help
	subFlags.Summary = sub.Summary

	if stop, err := subFlags.Boot(r.Flags.FlagSet.Args()); err != nil || stop() {
		return err
	}

	logx.Debugf("invoking %s in %s with %s",
		r.Flags.FlagSet.Args()[0], wd, r.Flags.FlagSet.Args())
	err = sub.Do(r.Flags.FlagSet.Args()[1:])
	return
}

func (r *RootCmd) Help() {
	fmt.Println("bar [OPTIONS] SUBCOMMAND [OPTIONS] ...\n")
}

func (r *RootCmd) Summary() {
	fmt.Fprintln(r.Stderr, "\nsubcommands\n")
	var names sort.StringSlice
	route := r.getRoute()
	for n, _ := range route {
		names = append(names, n)
	}
	names.Sort()
	for _, n := range names {
		fmt.Fprintf(r.Stderr, "  %s\n    \t%s\n", n,route[n](r.Base).Description())
	}
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
		"sneak":          NewSneakCmd,
	}
}
