package cmd
import (
	"flag"
	"io"
	"github.com/akaspin/bar/barctl/git"
	"github.com/akaspin/bar/barctl/transport"
	"github.com/akaspin/bar/shadow"
)

/*
Replace local shadows with downloaded BLOBs.

	$ barc down my/blobs

 */
type DownCmd struct {
	fs *flag.FlagSet
	endpoint string
	useGit bool

	git *git.Git
	transport *transport.TransportPool
	hasher *shadow.HasherPool
}

func (c *DownCmd) Bind(fs *flag.FlagSet, in io.Reader, out io.Writer) (err error) {
	c.fs = fs
	fs.StringVar(&c.endpoint, "endpoint", "http://localhost:3000/v1",
		"bard endpoint")
	fs.BoolVar(&c.useGit, "git", false, "use git infrastructure")
	return
}

func (c *DownCmd) Do() (err error) {
	if c.useGit {
		if c.git, err = git.NewGit(""); err != nil {
			return 
		}
	}

	return
}
