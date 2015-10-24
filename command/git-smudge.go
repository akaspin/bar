package command

import (
	"github.com/spf13/cobra"
	"github.com/tamtam-im/logx"
	"github.com/akaspin/bar/proto"
)

type GitSmudgeCmd struct  {
	*Environment
	*CommonOptions

}

func (c *GitSmudgeCmd) Init(cc *cobra.Command) {
	cc.Use = "smudge"
	cc.Short = "git smudge filter"
}

func (c *GitSmudgeCmd) Run(args ...string) (err error) {
	name := args[0]
	m, err := proto.NewFromAny(c.Stdin, c.ChunkSize)
	logx.Debugf("smudge manifest for %s (%s)", name, m.ID)
	err = m.Serialize(c.Stdout)
	return
}