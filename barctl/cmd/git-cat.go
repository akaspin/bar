package cmd
import (
	"flag"
	"github.com/akaspin/bar/shadow"
	"io"
	"fmt"
)

type GitCatCommand struct {
	fs *flag.FlagSet
	chunkSize int64
}

func (c *GitCatCommand) FS(fs *flag.FlagSet) {
	c.fs = fs
	fs.Int64Var(&c.chunkSize, "chunk-size", shadow.CHUNK_SIZE,
		"chunk size in bytes")
}

func (c *GitCatCommand) Do(in io.Reader, out, errOut io.Writer) (err error) {
	s, err := shadow.NewShadowFromFile(c.fs.Args()[0], false, c.chunkSize)
	if err != nil {
		return
	}
	fmt.Fprintf(out, "BAR-SHADOW-BLOB %x\n", s.ID)
	return
}

