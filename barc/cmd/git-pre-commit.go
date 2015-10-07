package cmd
import (
	"github.com/akaspin/bar/proto/manifest"
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

	endpoint string
	chunkSize int64
	pool int
}

func NewGitPreCommitCmd(s *BaseSubCommand) SubCommand {
	c := &GitPreCommitCmd{BaseSubCommand: s}
	c.FS.StringVar(&c.endpoint, "endpoint", "http://localhost:3000/v1",
		"bard endpoint")
	c.FS.Int64Var(&c.chunkSize, "chunk", manifest.CHUNK_SIZE, "preferred chunk size")
	c.FS.IntVar(&c.pool, "pool", 16, "pool size")
	return c
}

func (c *GitPreCommitCmd) Do() (err error) {
	return

}
