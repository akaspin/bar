package cmd
import (
	"io"
	"github.com/akaspin/bar/shadow"
	"flag"
	"fmt"
	"os"
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
		clean = barctl git-clean -endpoint=http://my.bar.server/v1 %f

	# .gitattributes
	my/blobs    filter=bar

STDIN is cat of file in working tree. Git will place STDOUT to
stage area.
*/
type GitCleanCommand struct {
	id bool
	silent bool
	fs *flag.FlagSet
	in io.Reader
	out io.Writer
}

func (c *GitCleanCommand) Bind(fs *flag.FlagSet, in io.Reader, out io.Writer) (err error) {
	c.fs = fs
	c.in, c.out = in, out

	fs.BoolVar(&c.id, "id", false, "print only id")
	fs.BoolVar(&c.silent, "silent", false, "supress warnings")

	return
}

func (c *GitCleanCommand) Do() (err error) {

	info, err := os.Stat(c.fs.Args()[0])
	if err != nil {
		return
	}

	var s *shadow.Shadow
	if s, err = shadow.New(c.in, info.Size()); err != nil {
		logx.Error(err)
		return
	}

	if s.IsFromShadow && c.silent {
		logx.Warning("warning %s is already shadow", c.fs.Args())
	}
	logx.Debugf("new shadow for %s %s", c.fs.Args()[0], s.ID)
	if c.id {
		fmt.Fprintf(c.out, "%s", s.ID)
	} else {
		err = s.Serialize(c.out)
	}

	return
}
