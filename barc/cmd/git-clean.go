package cmd
import (
	"github.com/akaspin/bar/proto/manifest"
	"fmt"
	"github.com/tamtam-im/logx"
	"github.com/akaspin/bar/barc/model"
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
	*BaseSubCommand

	id bool
	chunkSize int64
	pool int

	Model *model.Model
}

func NewGitCleanCommand(s *BaseSubCommand) SubCommand {
	c := &GitCleanCommand{BaseSubCommand: s}
	c.FS.BoolVar(&c.id, "id", false, "print only id")
	c.FS.Int64Var(&c.chunkSize, "chunk", manifest.CHUNK_SIZE, "preferred chunk size")
	c.FS.IntVar(&c.pool, "pool", 16, "pool size")
	return c
}

func (c *GitCleanCommand) Do() (err error) {
	c.Model, err = model.New(c.WD, true, c.chunkSize, c.pool)

	var name string
	if len(c.FS.Args()) > 0 {
		name = c.FS.Args()[0]
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
