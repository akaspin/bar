package command

import (
	"github.com/spf13/cobra"
	"github.com/tamtam-im/logx"
	"github.com/akaspin/bar/bar/model"
	"fmt"
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