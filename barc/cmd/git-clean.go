package cmd
import (
	"io"
	"github.com/akaspin/bar/shadow"
	"flag"
	"fmt"
	"github.com/tamtam-im/logx"
)

/*
Git clean. Used on `git add ...`. Equivalent of:

	$ cat my/file | barctl git-clean my/file

This command takes STDIN and filename and prints shadow
manifest to STDOUT.

To use with git register bar in git:

	# .git/config
	[filter "bar"]
		clean = barc git-clean -endpoint=http://my.bar.server/v1 %f

	# .gitattributes
	my/blobs    filter=bar

STDIN is cat of file in working tree. Git will place STDOUT to stage area.
Optional filename is always relative to git root.
*/
type GitCleanCommand struct {
	id bool
	chunkSize int64

	fs *flag.FlagSet
	in io.Reader
	out io.Writer
}

func (c *GitCleanCommand) Bind(wd string, fs *flag.FlagSet, in io.Reader, out io.Writer) (err error) {
	c.fs = fs
	c.in, c.out = in, out

	fs.BoolVar(&c.id, "id", false, "print only id")
	fs.Int64Var(&c.chunkSize, "chunk", shadow.CHUNK_SIZE, "preferred chunk size")
	return
}

func (c *GitCleanCommand) Do() (err error) {
	var name string
	if len(c.fs.Args()) > 0 {
		name = c.fs.Args()[0]
	}

	logx.Debugf("git-clean: %s", name)

	

	var s *shadow.Shadow
	if s, err = shadow.NewFromAny(c.in, c.chunkSize); err != nil {
		return
	}

	from := "BLOB"
	if s.IsFromShadow {
		from = "manifest"
	}
	logx.Debugf("shadow created from %s for %s %s", from, name, s.ID)

	if c.id {
		fmt.Fprintf(c.out, "%s", s.ID)
	} else {
		err = s.Serialize(c.out)
	}

	return
}
