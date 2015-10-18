package cmd
import (
	"fmt"
	"github.com/tamtam-im/logx"
	"github.com/akaspin/bar/bar/model"
	"flag"
)

/*
Git clean. Used on `git add ...`. Equivalent of:

	$ cat my/file | barctl git-clean my/file

This command takes STDIN and filename and prints shadow
manifest to STDOUT.

To use with git register bar in git:

	# .git/config
	[filter "bar"]
		clean = bar git-clean -endpoint=http://my.bar.server/v1 %f

	# .gitattributes
	my/blobs    filter=bar

STDIN is cat of file in working tree. Git will place STDOUT to stage area.
Optional filename is always relative to git root.
*/
type GitCleanCommand struct {
	*Base

	id bool

	Model *model.Model
}

func NewGitCleanCommand(s *Base) SubCommand {
	c := &GitCleanCommand{Base: s}
	return c
}

func (c *GitCleanCommand) Init(fs *flag.FlagSet) {
	fs.BoolVar(&c.id, "id", false, "print only id")
}

func (c *GitCleanCommand) Do(args []string) (err error) {
	c.Model, err = model.New(c.WD, true, c.ChunkSize, c.PoolSize)

	var name string
	if len(args) > 0 {
		name = args[0]
	}

	s, err := c.Model.GetManifest(name, c.Stdin)
	if err != nil {
		return
	}

	logx.Debugf("%s %s", name, s.ID)

	if c.id {
		fmt.Fprintf(c.Stdout, "%s", s.ID)
	} else {
		err = s.Serialize(c.Stdout)
	}
	return
}

func (c *GitCleanCommand) Description() string {
	return "git clean filter"
}