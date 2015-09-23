package cmd
import (
	"flag"
	"io"
	"github.com/akaspin/bar/barctl/git"
	"fmt"
	"strings"
)


// Git pre-commit hook. Used to upload all new/changed blobs
// to bard server:
//
// - Fails on uncommited bar-tracked BLOBs.
// - If working directory is clean - uploads BLOBs to bard.
type GitPreCommitCmd struct {
	endpoint string
}

func (c *GitPreCommitCmd) Bind(fs *flag.FlagSet, in io.Reader, out, errOut io.Writer) (err error) {
	return
}

func (c *GitPreCommitCmd) Do() (err error) {
	g, err := git.NewGit("")
	if err != nil {
		return
	}

	// Check dirty status
	dirty, err := g.DirtyFiles()
	if err != nil {
		return
	}
	dirty, err = g.FilterByDiff("bar", dirty...)
	if len(dirty) > 0 {
		err = fmt.Errorf("Dirty BLOBs in working tree. Run following command to add BLOBs:\n\n    git -C %s add %s\n",
			g.Root, strings.Join(dirty, " "))
	}

	// Collect BLOBs from diff

	return
}
