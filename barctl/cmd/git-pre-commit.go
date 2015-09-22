package cmd
import (
	"flag"
	"io"
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
	return
}
