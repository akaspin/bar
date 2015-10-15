package cmd
import (
	"github.com/akaspin/bar/barc/git"
	"github.com/akaspin/bar/barc/transport"
	"fmt"
	"github.com/tamtam-im/logx"
	"github.com/akaspin/bar/proto"
	"github.com/akaspin/bar/barc/model"
	"strings"
)

const hook  = `#!/bin/sh
# bar pre-commit hook
set -e

barc -log-level=%s git-pre-commit -http=%s -rpc=%s -chunk=%d -pool=%d
`

/*
Install bar for git infrastructure

	$ barc git-init -endpoint=http://my.bar.server/v1

This command installs git infrastructure to use with bar:

1. Adds bar filter to .git/config
2. Adds bar diff to .git/config
3. Adds pre-commit hook to .git/hooks
4. Adds git aliases for `barc up`, `barc down` and `barc ls`
*/
type GitInitCmd struct {
	*BaseSubCommand
	httpEndpoint string
	rpcEndpoints string
	log string
	clean bool

	git *git.Git
	transport *transport.Transport
}

func NewGitInitCmd(s *BaseSubCommand) SubCommand {
	c := &GitInitCmd{BaseSubCommand: s}
	c.FS.StringVar(&c.httpEndpoint, "http", "http://localhost:3000/v1",
		"bard http endpoint")
	c.FS.StringVar(&c.rpcEndpoints, "rpc", "http://localhost:3000/v1",
		"bard rpc endpoints separated by comma")
	c.FS.StringVar(&c.log, "log", "WARNING", "barc logging level")
	c.FS.BoolVar(&c.clean, "clean", false, "uninstall bar")
	return c
}

func (c *GitInitCmd) Do() (err error) {
	mod, err := model.New(c.WD, true, proto.CHUNK_SIZE, 16)
	if err != nil {
		return
	}

	c.git = mod.Git
	if c.clean {
		err = c.uninstall()
		return
	}

	c.transport = transport.NewTransport(mod, c.httpEndpoint, c.rpcEndpoints, 10)
	defer c.transport.Close()

	var opts proto.ServerInfo
	if opts, err = c.precheck(); err != nil {
		return
	}

	if err = c.git.SetHook("pre-commit",
		fmt.Sprintf(hook,
			c.log, c.httpEndpoint, strings.Join(opts.RPCEndpoints, ","),
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
//			"barc -log-level=%s git-diff -chunk=%d", c.log, info.ChunkSize),
		"filter.bar.clean": fmt.Sprintf(
			"barc -log-level=%s git-clean -chunk=%d -pool=%d %%f",
			c.log, info.ChunkSize, info.PoolSize),
		"filter.bar.smudge": fmt.Sprintf(
			"barc -log-level=%s git-smudge -chunk=%d -pool=%d %%f",
			c.log, info.ChunkSize, info.PoolSize),
		"alias.bar-squash": fmt.Sprintf(
			"!barc -log-level=%s up -squash -http=%s -rpc=%s -chunk=%d -pool=%d -git",
			c.log, info.HTTPEndpoint, rpc, info.ChunkSize, info.PoolSize),
		"alias.bar-up": fmt.Sprintf(
			"!barc -log-level=%s up -http=%s -rpc=%s -git -chunk=%d -pool=%d",
			c.log, info.HTTPEndpoint, rpc, info.ChunkSize, info.PoolSize),
		"alias.bar-down": fmt.Sprintf(
			"!barc -log-level=%s down -http=%s -rpc=%s -git -chunk=%d -pool=%d",
			c.log, info.HTTPEndpoint, rpc, info.ChunkSize, info.PoolSize),
		"alias.bar-ls": fmt.Sprintf(
			"!barc -log-level=%s ls -http=%s -rpc=%s -git -pool=%d",
			c.log, info.HTTPEndpoint, rpc, info.PoolSize),
		"alias.bar-export": fmt.Sprintf(
			"!barc -log-level=%s spec-export -upload -http=%s -rpc=%s -git -chunk=%d -pool=%d",
			c.log, info.HTTPEndpoint, rpc, info.ChunkSize, info.PoolSize),
		"alias.bar-import": fmt.Sprintf(
			"!barc -log-level=%s spec-import -http=%s -rpc=%s -git -chunk=%d -pool=%d",
			c.log, info.HTTPEndpoint, rpc, info.ChunkSize, info.PoolSize),
	}
}

// Prepare to install. Check endpoint, pre-commit hook
func (c *GitInitCmd) precheck() (res proto.ServerInfo, err error) {
	logx.Debugf("requesting endpoint %s", c.httpEndpoint)
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

