package cmd
import (
	"io"
	"flag"
)

/*
Git smudge. Used at `git checkout ...`. Equivalent of:

	$ cat staged/shadow | barctl git-smudge my/file

This command takes cat of shadow manifest as STDIN and FILE and
writes file contents to STDOUT.

To use with git register bar in git:

	# .git/config
	[filter "bar"]
		smudge = barctl git-smudge -endpoint=http://my.bar.server/v1 %f

	# .gitattributes
	my/blobs    filter=bar

STDIN is cat of file in git stage area. Git will place STDOUT to
working tree. If FILE is BLOB. git-smudge will check it ID and download
correct revision from bard. To disable this behaviour and replace changed
BLOBs with shadows use -shadow-changed option.
*/
type GitSmudgeCmd struct {
	endpoint string
	squashChanged bool

	out, errOut io.Writer
	in io.Reader
	fs *flag.FlagSet
}

func (c *GitSmudgeCmd) Bind(fs *flag.FlagSet, in io.Reader, out, errOut io.Writer) (err error) {
	c.fs = fs
	c.in, c.out, c.errOut = in, out, errOut

	fs.StringVar(&c.endpoint, "endpoint", "http://localhost:3000/v1",
		"bard endpoint")
	fs.BoolVar(&c.squashChanged, "shadow-changed", false,
		"replace changed blobs with shadows instead download from bard")
	return
}

func (c *GitSmudgeCmd) Do() (err error) {
	return
}
