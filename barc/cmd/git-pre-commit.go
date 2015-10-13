package cmd
import (
	"github.com/akaspin/bar/manifest"
	"github.com/akaspin/bar/barc/model"
	"fmt"
	"github.com/akaspin/bar/barc/transport"
)


/*
Git pre-commit hook. Used to upload all new/changed blobs
to bard server:

- Fails on uncommited bar-tracked BLOBs.
- If working directory is clean - uploads BLOBs to bard.

To use with git git-clean MUST be registered in git. Also
git pre-commit hook MUST be registered:

	$ cat > .git/hooks/pre-commit <<EOF
	#!/usr/bin/env sh
	set -e
	barctl git-pre-commit -endpoint=http://my.bar.server/v1
	EOF
	chmod +x .git/hooks/pre-commit
*/

type GitPreCommitCmd struct {
	*BaseSubCommand

	httpEndpoint string
	rpcEndpoints string
	chunkSize int64
	pool int

	model *model.Model
}

func NewGitPreCommitCmd(s *BaseSubCommand) SubCommand {
	c := &GitPreCommitCmd{BaseSubCommand: s}
	c.FS.StringVar(&c.httpEndpoint, "http", "http://localhost:3000/v1",
		"bard http endpoint")
	c.FS.StringVar(&c.rpcEndpoints, "rpc", "http://localhost:3000/v1",
		"bard rpc endpoints separated by comma")
	c.FS.Int64Var(&c.chunkSize, "chunk", manifest.CHUNK_SIZE, "preferred chunk size")
	c.FS.IntVar(&c.pool, "pool", 16, "pool size")
	return c
}

func (c *GitPreCommitCmd) Do() (err error) {
	if c.model, err = model.New(c.WD, true, c.chunkSize, c.pool); err != nil {
		return
	}

	isDirty, dirty, err := c.model.Check()
	if err != nil {
		return
	}
	if isDirty {
		err = fmt.Errorf("dirty files in working tree %s", dirty)
		return
	}

	feedR, err := c.model.Git.Diff()
	if err != nil {
		return
	}

	blobs, err := c.model.Git.ManifestsFromDiff(feedR)
	if err != nil {
		return
	}

	trans := transport.NewTransport(c.model, c.httpEndpoint, c.rpcEndpoints, c.pool)
	err = trans.Upload(blobs)

	return
}
