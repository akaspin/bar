package cmd
import (
	"github.com/akaspin/bar/proto"
	"github.com/tamtam-im/logx"
)

/*
Git smudge. Used at `git checkout ...`. Equivalent of:

	$ cat staged/shadow-or-BLOB | barctl git-smudge my/file

To use with git register bar in git:

	# .git/config
	[filter "bar"]
		smudge = barctl git-smudge -endpoint=http://my.bar.server/v1 %f

	# .gitattributes
	my/blobs    filter=bar

By default smudge just parse manifest from STDIN and pass to STDOUT. If STDIN
is BLOB - it will be uploaded to bard.
*/
type GitSmudgeCmd struct {
	*BaseSubCommand

	chunkSize int64
	maxConn int
}

func NewGitSmudgeCmd(s *BaseSubCommand) SubCommand {
	c := &GitSmudgeCmd{BaseSubCommand: s}

	c.FS.Int64Var(&c.chunkSize, "chunk", proto.CHUNK_SIZE, "preferred chunk size")
	c.FS.IntVar(&c.maxConn, "pool", 16, "pool size")
	return c
}

func (c *GitSmudgeCmd) Do() (err error) {
	name := c.FS.Args()[0]
//	r, isManifest, err := shadow.Peek(c.in)
//	if isManifest {
//
//	}

	m, err := proto.NewFromAny(c.Stdin, c.chunkSize)
	logx.Debugf("smudge manifest for %s", name, m.ID)
	err = m.Serialize(c.Stdout)
	return
}
