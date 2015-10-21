package git
import (
	"github.com/akaspin/bar/proto"
	"html/template"
	"strings"
	"bytes"
	"fmt"
	"github.com/tamtam-im/logx"
)

const preCommitHook  = `#!/bin/sh
# bar pre-commit hook
set -e

bar git pre-commit --log-level={{.LogLevel}} --endpoint={{Join .Info.RPCEndpoints ","}} --chunk={{.Info.ChunkSize}} --pool={{.Info.PoolSize}} --buffer={{.Info.BufferSize}}
`

// Bar git configuration
type Config struct {
	proto.ServerInfo
	*Git
}

func NewConfig(info proto.ServerInfo, g *Git) *Config {
	return &Config{info, g}
}

func (c *Config) Install(logLevel string) (err error) {
	// check hook
	_, hookErr := c.Git.GetHook("pre-commit")
	if hookErr == nil {
		err = fmt.Errorf("pre-commit hook already exists")
	}

	h, err := c.runInfoTpl(preCommitHook, logLevel, c.ServerInfo)
	if err != nil {
		return
	}

	if err = c.Git.SetHook("pre-commit", h); err != nil {
		return
	}

	for k, v := range c.getConfigLines() {
		var v1 string
		if v1, err = c.runInfoTpl(v, logLevel, c.ServerInfo); err != nil {
			return
		}
		if err = c.Git.SetConfig(k, v1); err != nil {
			return
		}
		logx.Debugf("set %s = %s", k, v1)
	}

	return
}

func (c *Config) Uninstall() (err error) {
	c.Git.CleanHook("pre-commit")
	for k, _ := range c.getConfigLines() {
		err = c.Git.UnsetConfig(k)
		logx.Debugf("unset %s", k)
	}
	return
}

func (c *Config) getConfigLines() map[string]string {
	return map[string]string{
		// filters
		"filter.bar.clean": `bar git clean --log-level={{.LogLevel}} --endpoint={{Join .Info.RPCEndpoints ","}} --chunk={{.Info.ChunkSize}} --pool={{.Info.PoolSize}} %f`,
		"filter.bar.smudge": `bar git clean --log-level={{.LogLevel}} --endpoint={{Join .Info.RPCEndpoints ","}} --chunk={{.Info.ChunkSize}} --pool={{.Info.PoolSize}} %f`,

		// basic
		"alias.bar-up": `!bar up --git --log-level={{.LogLevel}} --endpoint={{Join .Info.RPCEndpoints ","}} --chunk={{.Info.ChunkSize}} --pool={{.Info.PoolSize}} --buffer={{.Info.BufferSize}}`,
		"alias.bar-down": `!bar down --git --log-level={{.LogLevel}} --endpoint={{Join .Info.RPCEndpoints ","}} --chunk={{.Info.ChunkSize}} --pool={{.Info.PoolSize}} --buffer={{.Info.BufferSize}}`,
		"alias.bar-squash": `!bar up --git --squash --log-level={{.LogLevel}} --endpoint={{Join .Info.RPCEndpoints ","}} --chunk={{.Info.ChunkSize}} --pool={{.Info.PoolSize}} --buffer={{.Info.BufferSize}}`,
		"alias.bar-ls": `!bar ls --git --log-level={{.LogLevel}} --endpoint={{Join .Info.RPCEndpoints ","}} --chunk={{.Info.ChunkSize}} --pool={{.Info.PoolSize}} --buffer={{.Info.BufferSize}}`,

		// spec
		"alias.bar-spec-export": `!bar spec export --git --upload --log-level={{.LogLevel}} --endpoint={{Join .Info.RPCEndpoints ","}} --chunk={{.Info.ChunkSize}} --pool={{.Info.PoolSize}} --buffer={{.Info.BufferSize}}`,
		"alias.bar-spec-import": `!bar spec import --git --squash --log-level={{.LogLevel}} --endpoint={{Join .Info.RPCEndpoints ","}} --chunk={{.Info.ChunkSize}} --pool={{.Info.PoolSize}} --buffer={{.Info.BufferSize}}`,
		"alias.bar-spec-ls": `!bar spec import --git --log-level={{.LogLevel}} --endpoint={{Join .Info.RPCEndpoints ","}} --chunk={{.Info.ChunkSize}} --pool={{.Info.PoolSize}} --buffer={{.Info.BufferSize}}`,
	}
}

func (c *Config) runInfoTpl(what string, logLevel string, info proto.ServerInfo) (res string, err error) {
	t, err := template.New("tpl").Funcs(
		template.FuncMap{"Join": strings.Join}).Parse(what)
	if err != nil {
		return
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, map[string]interface{}{
		"LogLevel": logLevel,
		"Info": info,
	}); err != nil {
		return
	}
	res = string(buf.Bytes())
	return
}