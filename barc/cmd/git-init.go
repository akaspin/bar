package cmd
import (
	"github.com/akaspin/bar/barc/git"
	"github.com/akaspin/bar/barc/transport"
	"fmt"
	"github.com/tamtam-im/logx"
	"github.com/akaspin/bar/proto"
)

const hook  = `#!/bin/sh
# bar pre-commit hook
set -e

barc -log-level=%s git-pre-commit -endpoint=%s -chunk=%d -pool=%d
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
	endpoint string
	log string
	clean bool

	git *git.Git
	transport *transport.Transport
}

func NewGitInitCmd(s *BaseSubCommand) SubCommand {
	c := &GitInitCmd{BaseSubCommand: s}
	c.FS.StringVar(&c.endpoint, "endpoint", "http://localhost:3000/v1",
		"bard endpoint")
	c.FS.StringVar(&c.log, "log", "WARNING", "barc logging level")
	c.FS.BoolVar(&c.clean, "clean", false, "uninstall bar")
	return c
}

func (c *GitInitCmd) Do() (err error) {
	c.git, err = git.NewGit(c.WD)
	if c.clean {
		err = c.uninstall()
		return
	}

	c.transport = transport.NewTransport(c.WD, c.endpoint, 10)
	defer c.transport.Close()

	var opts proto.Info
	if opts, err = c.precheck(); err != nil {
		return
	}

	if err = c.git.SetHook("pre-commit",
		fmt.Sprintf(hook,
			c.log, c.endpoint, opts.ChunkSize, opts.MaxConn)); err != nil {
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

func (c *GitInitCmd) configVals(info proto.Info) map[string]string {
	return map[string]string{
//		"diff.bar.command": fmt.Sprintf(
//			"barc -log-level=%s git-diff -chunk=%d", c.log, info.ChunkSize),
		"filter.bar.clean": fmt.Sprintf(
			"barc -log-level=%s git-clean -chunk=%d -pool=%s %%f",
			c.log, info.ChunkSize, info.MaxConn),
		"filter.bar.smudge": fmt.Sprintf(
			"barc -log-level=%s git-smudge -endpoint=%s -chunk=%d -pool=%d %%f",
			c.log, c.endpoint, info.ChunkSize, info.MaxConn),
//		"alias.bar-squash": fmt.Sprintf(
//			"!barc -log-level=%s up -squash -endpoint=%s -chunk=%d -pool=%d -git",
//			c.log, c.endpoint, info.ChunkSize, info.MaxConn),
		"alias.bar-up": fmt.Sprintf(
			"!barc -log-level=%s up -endpoint=%s -git -chunk=%d -pool=%d",
			c.log, c.endpoint, info.ChunkSize, info.MaxConn),
//		"alias.bar-down": fmt.Sprintf(
//			"!barc -log-level=%s down -endpoint=%s -git -chunk=%d -pool=%d",
//			c.log, c.endpoint, info.ChunkSize, info.MaxConn),
//		"alias.bar-ls": fmt.Sprintf(
//			"!barc -log-level=%s ls -endpoint=%s -git -pool=%d",
//			c.log, c.endpoint, info.MaxConn),
	}
}

// Prepare to install. Check endpoint, pre-commit hook
func (c *GitInitCmd) precheck() (res proto.Info, err error) {
	if res, err = c.transport.Ping(); err != nil {
		return
	}

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
	for k, _ := range c.configVals(proto.Info{}) {
		logx.Debugf("removing %s from git config", k)
		c.git.UnsetConfig(k)
	}

	return
}

