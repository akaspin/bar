package command
import (
	"github.com/spf13/cobra"
	"github.com/tamtam-im/logx"
	"github.com/akaspin/bar/bar/model"
	"github.com/akaspin/bar/bar/transport"
	"github.com/akaspin/bar/bar/git"
)

type GitInstallCmd struct  {
	*Environment
	*CommonOptions

	// Installable logging lenel
	Log string

	// Diversion camp location
//	DivertDir string
}

func (c *GitInstallCmd) Init(cc *cobra.Command) {
	cc.Use = "install"
	cc.Short = "install bar support into git repo"

	cc.Flags().StringVarP(&c.Log, "log", "", logx.DEBUG,
		"installable logging level")

	cc.Flags()
}

func (c *GitInstallCmd) Run(args ...string) (err error) {
	var mod *model.Model
	if mod, err = model.New(c.WD, true, c.ChunkSize, c.PoolSize); err != nil {
		return
	}
	defer mod.Close()

	trans := transport.NewTransport(mod, "", c.Endpoint, c.PoolSize)
	defer trans.Close()

	info, err := trans.ServerInfo()
	if err != nil {
		return
	}

	config := git.NewConfig(info, mod.Git)
	err = config.Install(c.Log)

	return
}

