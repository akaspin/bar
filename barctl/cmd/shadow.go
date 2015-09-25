package cmd
import (
	"flag"
	"io"
)


/*
This command upload BLOBs to bard and replaces
them with shadows.

	$ barctl matter my/blobs my/blobs/glob*


*/
type ShadowCmd struct {
	git bool
	endpoint string
	poolSize int

	fs *flag.FlagSet
}

func (c *ShadowCmd) Bind(fs *flag.FlagSet, in io.Reader, out io.Writer) (err error) {
	c.fs = fs

	fs.StringVar(&c.endpoint, "endpoint", "http://localhost:3000/v1",
		"bard endpoint")
	fs.BoolVar(&c.git, "git", false, "use git attributes to restrict globals")
	fs.IntVar(&c.poolSize, "pool", 16, "pool size")

	return
}
