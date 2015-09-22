package cmd
import (
	"flag"
	"github.com/akaspin/bar/shadow"
	"io"
	"fmt"
	"os"
	"github.com/akaspin/bar/barctl/git"
	"strings"
)

// Cat for commit diff.
//
//    $ git diff-files <file>  # fail on any output
//    $ git ls-files --cached -s --full-name <file>
//    100644 0972a66281ba8cee7bb6ad3ad322a9afe6830338 0	fixtures/roygbiv.jpg
//           ----------------------------------------
//           Find OID
//    $ git cat-file -p 9ccd85cc5461042dbc2db1ea43ab81558c7b1710
//    ...
//    Get BLOB ID from staging area
//
type GitCatCommand struct {
	fs *flag.FlagSet
	out, errOut io.Writer
	strict bool
	chunkSize int64
}

func (c *GitCatCommand) Bind(fs *flag.FlagSet, in io.Reader, out, errOut io.Writer) (err error) {
	c.fs = fs
	c.out = out
	c.errOut = errOut
	return
}

func (c *GitCatCommand) Do() (err error) {
	n := c.fs.Args()[0]
	outsider := strings.HasPrefix(n, "/")
	info, err := os.Stat(c.fs.Args()[0])
	if err != nil {
		return
	}

	gitRoot, err := git.GetGitTop()
	if err != nil {
		return
	}

	var r io.Reader
	if outsider {
		// outsider
		r, err = os.Open(n)
		if err != nil {
			return
		}
		defer r.(*os.File).Close()
	} else {
		var oid string
		oid, err = git.GetFileOID(gitRoot, n)
		if err != nil {
			return
		}
		r, err = git.CatFile(oid)
		if err != nil {
			return
		}
	}

	s, err := shadow.New(r, info.Size())
	if err != nil {
		return
	}

	fmt.Fprintf(c.out, "BAR-SHADOW-BLOB %s\n", s.ID)
	return
}

