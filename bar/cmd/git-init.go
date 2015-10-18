package cmd
import (
	"github.com/akaspin/bar/bar/git"
	"github.com/akaspin/bar/bar/transport"
	"fmt"
	"github.com/tamtam-im/logx"
	"github.com/akaspin/bar/proto"
	"github.com/akaspin/bar/bar/model"
	"strings"
	"flag"
)

const hook  = `#!/bin/sh
# bar pre-commit hook
set -e

bar -log-level=%s -endpoint=%s -chunk=%d -pool=%d git-pre-commit
`

/*
Install bar for git infrastructure

	$ bar git-init -endpoint=http://my.bar.server/v1

This command installs git infrastructure to use with bar:

1. Adds bar filter to .git/config
2. Adds bar diff to .git/config
3. Adds pre-commit hook to .git/hooks
4. Adds git aliases for `bar up`, `bar down` and `bar ls`
*/
type GitInitCmd struct {
	*Base
	log string
	clean bool

	git *git.Git
	transport *transport.Transport
}

func NewGitInitCmd(s *Base) SubCommand {
	c := &GitInitCmd{Base: s}
	return c
}

func (c *GitInitCmd) Init(fs *flag.FlagSet) {
	fs.StringVar(&c.log, "log", logx.INFO, "bar logging level")
	fs.BoolVar(&c.clean, "clean", false, "uninstall bar")
}

func (c *GitInitCmd) Description() string {
	return "install bar to git repo"
}

func (c *GitInitCmd) Do(args []string) (err error) {
	mod, err := model.New(c.WD, true, proto.CHUNK_SIZE, 16)
	if err != nil {
		return
	}

	c.git = mod.Git
	if c.clean {
		err = c.uninstall()
		return
	}

	c.transport = transport.NewTransport(mod, "", c.BaseOptions.Endpoints, 10)
	defer c.transport.Close()

	var opts proto.ServerInfo
	if opts, err = c.precheck(); err != nil {
		return
	}

	if err = c.git.SetHook("pre-commit",
		fmt.Sprintf(hook,
			c.log, "", strings.Join(opts.RPCEndpoints, ","),
			opts.ChunkSize, opts.PoolSize)); err != nil {
		return
	}
	logx.Infof("pre-commit hook installed to %s",
		c.git.Root + ".git/hooks/pre-commit")

	for k, v := range c.configVals(opts) {
		c.git.SetConfig(k, v)
		logx.Debugf("setting git config %s %s", k, v)
	}

	return
}

func (c *GitInitCmd) configVals(info proto.ServerInfo) map[string]string {
	rpc := strings.Join(info.RPCEndpoints, ",")
	return map[string]string{
//		"diff.bar.command": fmt.Sprintf(
//			"bar -log-level=%s git-diff -chunk=%d", c.log, info.ChunkSize),
		"filter.bar.clean": fmt.Sprintf(
			"bar -log-level=%s -chunk=%d -pool=%d git-clean %%f",
			c.log, info.ChunkSize, info.PoolSize),
		"filter.bar.smudge": fmt.Sprintf(
			"bar -log-level=%s -chunk=%d -pool=%d git-smudge %%f",
			c.log, info.ChunkSize, info.PoolSize),
		"alias.bar-squash": fmt.Sprintf(
			"!bar -log-level=%s -endpoint=%s -chunk=%d -pool=%d up -squash -git",
			c.log, rpc, info.ChunkSize, info.PoolSize),
		"alias.bar-up": fmt.Sprintf(
			"!bar -log-level=%s -endpoint=%s -chunk=%d -pool=%d up -git",
			c.log, rpc, info.ChunkSize, info.PoolSize),
		"alias.bar-down": fmt.Sprintf(
			"!bar -log-level=%s -endpoint=%s -chunk=%d -pool=%d down -git",
			c.log, rpc, info.ChunkSize, info.PoolSize),
		"alias.bar-ls": fmt.Sprintf(
			"!bar -log-level=%s -endpoint=%s -chunk=%d -pool=%d ls -git",
			c.log, rpc, info.PoolSize),
		"alias.bar-export": fmt.Sprintf(
			"!bar -log-level=%s -endpoint=%s -chunk=%d -pool=%d spec-export -upload -git",
			c.log, rpc, info.ChunkSize, info.PoolSize),
		"alias.bar-import": fmt.Sprintf(
			"!bar -log-level=%s -endpoint=%s -chunk=%d -pool=%d spec-import -git -squash",
			c.log, rpc, info.ChunkSize, info.PoolSize),
		"alias.bar-spec-ls": fmt.Sprintf(
			"!bar -log-level=%s -endpoint=%s -chunk=%d -pool=%d spec-import -git -squash",
			c.log, rpc, info.ChunkSize, info.PoolSize),
	}
}

// Prepare to install. Check endpoint, pre-commit hook
func (c *GitInitCmd) precheck() (res proto.ServerInfo, err error) {
	if res, err = c.transport.ServerInfo(); err != nil {
		return
	}

	logx.Debug(res)

	// Check for hook - fail if exists
	_, hookErr := c.git.GetHook("pre-commit")
	if hookErr == nil {
		err = fmt.Errorf("pre-commit hook already exists")
	}
	return
}

func (c *GitInitCmd) uninstall() (err error) {
	logx.Debug("removing pre-commit hook")
	c.git.CleanHook("pre-commit")
	for k, _ := range c.configVals(proto.ServerInfo{}) {
		logx.Debugf("removing %s from git config", k)
		c.git.UnsetConfig(k)
	}

	return
}

