package command

import (
	"github.com/spf13/cobra"
	"github.com/tamtam-im/logx"
)

type GitUninstallCmd struct  {
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
	logx.Debug(c)
	return
}