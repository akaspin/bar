package command
import (
	"github.com/spf13/cobra"
	"github.com/tamtam-im/logx"
)

type GitInstallCmd struct  {
	*Environment
	*CommonOptions

	// Installable logging lenel
	Log string
}

func (c *GitInstallCmd) Init(cc *cobra.Command) {
	cc.Use = "install"
	cc.Short = "install bar support into git repo"
	cc.Flags().StringVarP(&c.Log, "log", "", logx.DEBUG,
		"installable logging level")
}

func (c *GitInstallCmd) Run(args ...string) (err error) {
	logx.Error(c, c.CommonOptions)
	return
}