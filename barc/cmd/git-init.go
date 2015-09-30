package cmd
import (
	"github.com/akaspin/bar/barc/git"
	"flag"
	"io"
	"github.com/akaspin/bar/barc/transport"
	"net/url"
	"time"
	"fmt"
	"github.com/tamtam-im/logx"
)

const hook  = `#!/bin/sh
# bar pre-commit hook
set -e

barc -log-level=%s git-pre-commit -endpoint=%s
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
	endpoint string
	log string
	clean bool

	git *git.Git
	transport *transport.TransportPool
}

func (c *GitInitCmd) Bind(wd string, fs *flag.FlagSet, in io.Reader, out io.Writer) (err error) {
	fs.StringVar(&c.endpoint, "endpoint", "http://localhost:3000/v1",
		"bard endpoint")
	fs.StringVar(&c.log, "log", "WARNING", "barc logging level")
	fs.BoolVar(&c.clean, "clean", false, "uninstall bar")
	c.git, err = git.NewGit("")
	return
}

func (c *GitInitCmd) Do() (err error) {
	if c.clean {
		err = c.uninstall()
		return
	}

	// init transport
	u, err := url.Parse(c.endpoint)
	if err != nil {
		return
	}
	c.transport = transport.NewTransportPool(u, 10, time.Minute)

	if err = c.precheck(); err != nil {
		return
	}

	if err = c.git.SetHook("pre-commit",
		fmt.Sprintf(hook, c.log, c.endpoint)); err != nil {
		return
	}
	logx.Infof("pre-commit hook installed to %s",
		c.git.Root + ".git/hooks/pre-commit")

	for k, v := range c.configVals() {
		c.git.SetConfig(k, v)
	}

	return
}

func (c *GitInitCmd) configVals() map[string]string {
	return map[string]string{
		"diff.bar.textconv": fmt.Sprintf(
			"barc -log-level=%s git-cat", c.log),
		"filter.bar.clean": fmt.Sprintf(
			"barc -log-level=%s git-clean %%f", c.log),
		"filter.bar.smudge": fmt.Sprintf(
			"barc -log-level=%s git-smudge -endpoint=%s %%f", c.log, c.endpoint),
		"alias.bar-squash": fmt.Sprintf(
			"!barc -log-level=%s up -endpoint=%s -git", c.log, c.endpoint),
		"alias.bar-up": fmt.Sprintf(
			"!barc -log-level=%s up -noop -endpoint=%s -git", c.log, c.endpoint),
		"alias.bar-down": fmt.Sprintf(
			"!barc -log-level=%s down -endpoint=%s -git", c.log, c.endpoint),
		"alias.bar-ls": fmt.Sprintf(
			"!barc -log-level=%s ls -endpoint=%s -git", c.log, c.endpoint),
	}
}

// Prepare to install. Check endpoint, pre-commit hook
func (c *GitInitCmd) precheck() (err error) {
	tr, err := c.transport.Take()
	if err != nil {
		return
	}
	defer c.transport.Release(tr)

	if err = tr.Ping(); err != nil {
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
	c.git.CleanHook("pre-commit")
	for k, _ := range c.configVals() {
		c.git.UnsetConfig(k)
	}

	return
}

