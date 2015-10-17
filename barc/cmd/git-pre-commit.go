package cmd
import (
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
	*Base

	model *model.Model
}

func NewGitPreCommitCmd(s *Base) SubCommand {
	c := &GitPreCommitCmd{Base: s}
	return c
}

func (c *GitPreCommitCmd) Do(args []string) (err error) {
	if c.model, err = model.New(c.WD, true, c.ChunkSize, c.PoolSize); err != nil {
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

	trans := transport.NewTransport(c.model, "", c.Endpoints, c.PoolSize)
	err = trans.Upload(blobs)

	return
}
