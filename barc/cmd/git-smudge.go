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
	*Base

}

func NewGitSmudgeCmd(s *Base) SubCommand {
	c := &GitSmudgeCmd{Base: s}
	return c
}

func (c *GitSmudgeCmd) Do(args []string) (err error) {
	name := args[0]

	m, err := proto.NewFromAny(c.Stdin, c.ChunkSize)
	logx.Debugf("smudge manifest for %s", name, m.ID)
	err = m.Serialize(c.Stdout)
	return
}
