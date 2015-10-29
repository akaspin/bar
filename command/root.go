package command

import (
	"github.com/akaspin/bar/proto"
	"github.com/spf13/cobra"
	"github.com/tamtam-im/logx"
	"io"
	"os"
	"path/filepath"
)

type RootCmd struct {
	*Environment
	*CommonOptions
}

func (c *RootCmd) Init(cc *cobra.Command) {
	cc.Use = "bar"
	pf := cc.PersistentFlags()
	pf.StringVarP(&c.LoggingLevel, "log-level", "", logx.INFO,
		"logging level")
	pf.StringVarP(&c.WD, "work-dir", "C", "", "work directory")

	pf.StringVarP(&c.Endpoint, "endpoint", "", ":3001",
		"bar server endpoint")
	pf.Int64VarP(&c.ChunkSize, "chunk", "", proto.CHUNK_SIZE,
		"preferred chunk size")
	pf.IntVarP(&c.PoolSize, "pool", "", 16, "pool size")
	pf.IntVarP(&c.BufferSize, "buffer", "", 1024*1024*8,
		"thrift buffer size per connection")

	cc.PersistentPreRunE = func(cc *cobra.Command, args []string) (err error) {
		logx.SetOutput(c.Stderr)
		logx.SetLevel(c.LoggingLevel)
		wd, err := os.Getwd()
		if err != nil {
			return
		}
		c.WD = filepath.Clean(filepath.Join(wd, c.WD))
		logx.Debugf("working in %s : %s", c.WD, args)
		return
	}
	return
}

func (c *RootCmd) Run(args ...string) (err error) { return }

func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) (err error) {
	env := &Environment{stdin, stdout, stderr}
	cOpts := &CommonOptions{}
	serverOpts := &ServerOptions{}

	root := Attach(
		&RootCmd{Environment: env, CommonOptions: cOpts}, env,
		Attach(&LsCmd{Environment: env, CommonOptions: cOpts}, env),
		Attach(&UpCmd{Environment: env, CommonOptions: cOpts}, env),
		Attach(&DownCmd{Environment: env, CommonOptions: cOpts}, env),
		Attach(&GitCmd{}, env,
			Attach(&GitInstallCmd{Environment: env, CommonOptions: cOpts}, env),
			Attach(&GitUninstallCmd{Environment: env, CommonOptions: cOpts}, env),
			Attach(&GitCleanCmd{Environment: env, CommonOptions: cOpts}, env),
			Attach(&GitSmudgeCmd{Environment: env, CommonOptions: cOpts}, env),
			Attach(&GitPreCommitCmd{Environment: env, CommonOptions: cOpts}, env),
		),
		Attach(&GitDivertRootCmd{}, env,
			Attach(&GitDivertStatusCmd{Environment: env, CommonOptions: cOpts}, env),
			Attach(&GitDivertBeginCmd{Environment: env, CommonOptions: cOpts}, env),
			Attach(&GitDivertFinishCmd{Environment: env, CommonOptions: cOpts}, env),
			Attach(&GitDivertAbortCmd{Environment: env, CommonOptions: cOpts}, env),
			Attach(&GitDivertPushCmd{Environment: env, CommonOptions: cOpts}, env),
		),
		Attach(&SpecRootCmd{}, env,
			Attach(&SpecExportCmd{Environment: env, CommonOptions: cOpts}, env),
			Attach(&SpecImportCmd{Environment: env, CommonOptions: cOpts}, env),
		),
		Attach(&VersionCmd{Environment: env}, env),
		Attach(&PingCmd{Environment: env, CommonOptions: cOpts}, env),
		Attach(&ServerCmd{}, env,
			Attach(&ServerRunCmd{
				Environment:   env,
				CommonOptions: cOpts,
				ServerOptions: serverOpts,
			}, env),
		),
	)
	root.SetArgs(args)
	err = root.Execute()
	return
}
