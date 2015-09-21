package cmd
import (
	"flag"
	"github.com/akaspin/bar/shadow"
	"io"
	"fmt"
"github.com/akaspin/bar/fixtures"
)

type GitCatCommand struct {
	fs *flag.FlagSet
	out io.Writer

	chunkSize int64
}

func (c *GitCatCommand) Bind(fs *flag.FlagSet, in io.Reader, out, errOut io.Writer) (err error) {
	c.fs = fs
	c.out = out
	fs.Int64Var(&c.chunkSize, "chunk-size", shadow.CHUNK_SIZE,
		"chunk size in bytes")
	return
}

func (c *GitCatCommand) Do() (err error) {
	s, err := fixtures.NewShadowFromFile(c.fs.Args()[0], false, c.chunkSize)
	if err != nil {
		return
	}
	fmt.Fprintf(c.out, "BAR-SHADOW-BLOB %s\n", s.ID)
	return
}

