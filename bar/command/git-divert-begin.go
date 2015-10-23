package command
import (
	"github.com/spf13/cobra"
	"github.com/tamtam-im/logx"
	"fmt"
	"github.com/akaspin/bar/bar/model"
	"github.com/akaspin/bar/bar/git"
)

type GitDivertBeginCmd struct  {
	*Environment
	*CommonOptions

}

func (c *GitDivertBeginCmd) Init(cc *cobra.Command) {
	cc.Use = "begin BRANCH [# TREE-ISH ...]"
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
	if err = divert.Begin(branch, names...); err == nil {
		return
	}

	return
}