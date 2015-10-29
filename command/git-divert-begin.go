package command

import (
	"fmt"
	"github.com/akaspin/bar/client/git"
	"github.com/akaspin/bar/client/model"
	"github.com/spf13/cobra"
	"github.com/tamtam-im/logx"
)

type GitDivertBeginCmd struct {
	*Environment
	*CommonOptions
}

func (c *GitDivertBeginCmd) Init(cc *cobra.Command) {
	cc.Use = "begin branch [# tree-ish ...]"
	cc.Short = "begin covert op"
}

func (c *GitDivertBeginCmd) Run(args ...string) (err error) {
	logx.Debugf("beginning covert op %s", args)

	if len(args) == 0 {
		err = fmt.Errorf("no branch")
		return
	}

	branch := args[0]
	names := args[1:]

	mod, err := model.New(c.WD, true, c.ChunkSize, c.PoolSize)
	if err != nil {
		return
	}
	divert := git.NewDivert(mod.Git)
	var spec git.DivertSpec
	if spec, err = divert.PrepareBegin(branch, names...); err != nil {
		return
	}

	if err = divert.Begin(spec); err == nil {
		return
	}

	return
}
