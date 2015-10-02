package cmd
import (
	"io"
	"flag"
	"github.com/akaspin/bar/shadow"
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
	endpoint string
	chunkSize int64
	maxConn int

	fs *flag.FlagSet
	in io.Reader
	out io.Writer

}

func (c *GitSmudgeCmd) Bind(wd string, fs *flag.FlagSet, in io.Reader, out io.Writer) (err error) {
	c.fs = fs
	c.in, c.out = in, out

	fs.StringVar(&c.endpoint, "endpoint", "http://localhost:3000/v1",
		"bard endpoint")
	fs.Int64Var(&c.chunkSize, "chunk", shadow.CHUNK_SIZE, "preferred chunk size")
	fs.IntVar(&c.maxConn, "pool", 16, "pool size")
	return
}

func (c *GitSmudgeCmd) Do() (err error) {
	name := c.fs.Args()[0]
//	r, isManifest, err := shadow.Peek(c.in)
//	if isManifest {
//
//	}

	m, err := shadow.NewFromAny(c.in, c.chunkSize)
	logx.Debugf("smudge manifest for %s (ID: %s from-shadow: %t)",
		name, m.ID, m.IsFromShadow)
	err = m.Serialize(c.out)
	return
}
