package command

import (
	"github.com/akaspin/bar/client/git"
	"github.com/akaspin/bar/client/model"
	"github.com/akaspin/bar/proto"
	"github.com/spf13/cobra"
)

type GitUninstallCmd struct {
	*Environment
	*CommonOptions

	// Installable logging lenel
	Log string
}

func (c *GitUninstallCmd) Init(cc *cobra.Command) {
	cc.Use = "uninstall"
	cc.Short = "remove bar support from git repo"
}

func (c *GitUninstallCmd) Run(args ...string) (err error) {
	var mod *model.Model

	if mod, err = model.New(c.WD, true, c.ChunkSize, c.PoolSize); err != nil {
		return
	}

	config := git.NewConfig(proto.ServerInfo{}, mod.Git)
	err = config.Uninstall()

	return
}
