package command

import (
	"fmt"
	"github.com/akaspin/bar/client/model"
	"github.com/spf13/cobra"
	"github.com/tamtam-im/logx"
)

type GitDivertPushCmd struct {
	*Environment
	*CommonOptions
}

func (c *GitDivertPushCmd) Init(cc *cobra.Command) {
	cc.Use = "push [upstream] branch"
	cc.Short = "push specific branch"
}

func (c *GitDivertPushCmd) Run(args ...string) (err error) {
	mod, err := model.New(c.WD, true, c.ChunkSize, c.PoolSize)
	if err != nil {
		return
	}

	var upstream, branch string
	if len(args) == 0 {
		err = fmt.Errorf("no upstream and/or branch provided")
		return
	}
	if len(args) == 1 {
		upstream = "origin"
		branch = args[0]
	} else {
		upstream = args[0]
		branch = args[1]
	}

	// checks
	current, branches, err := mod.Git.GetBranches()
	if err != nil {
		return
	}
	if branch == current {
		err = fmt.Errorf("cannot push current branch. use `git push ...`")
		return
	}
	var exists bool
	for _, i := range branches {
		if branch == i {
			exists = true
			break
		}
	}
	if !exists {
		err = fmt.Errorf("branch %s is not exists")
		return
	}

	if err = mod.Git.Push(upstream, branch); err != nil {
		return
	}
	logx.Debugf("%s/%s pushed", upstream, branch)
	return
}
