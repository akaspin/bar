package command

import (
	"github.com/spf13/cobra"
	"github.com/tamtam-im/logx"
	"github.com/akaspin/bar/bar/model"
	"fmt"
	"github.com/akaspin/bar/bar/git"
)

type GitCleanCmd struct  {
	*Environment
	*CommonOptions

	// Installable logging lenel
	Id bool
}

func (c *GitCleanCmd) Init(cc *cobra.Command) {
	cc.Use = "clean"
	cc.Short = "git clean filter"
	cc.Flags().BoolVarP(&c.Id, "id", "", false, "print generated id instead manifest")
}

func (c *GitCleanCmd) Run(args ...string) (err error) {
	mod, err := model.New(c.WD, true, c.ChunkSize, c.PoolSize)

	var name string
	if len(args) > 0 {
		name = args[0]
	}

	// check divert
	divert := git.NewDivert(mod.Git)
	isInProgress, err := divert.IsInProgress()
	if err != nil {
		return
	}
	if isInProgress {
		var spec git.DivertSpec
		if spec, err = divert.ReadSpec(); err != nil {
			return
		}
		var exists bool
		for _, n := range spec.TargetFiles {
			if n == name {
				exists = true
				break
			}
		}
		if !exists {
			err = fmt.Errorf("wan't clean non-target file %s while divert in progress", name)
			return
		}
	}

	s, err := mod.GetManifest(name, c.Stdin)
	if err != nil {
		return
	}

	logx.Debugf("%s %s", name, s.ID)

	if c.Id {
		fmt.Fprintf(c.Stdout, "%s", s.ID)
	} else {
		err = s.Serialize(c.Stdout)
	}
	return
}