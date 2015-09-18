package cmd
import (
	"io"
	"github.com/akaspin/bar/shadow"
	"flag"
	"fmt"
)

type GitCleanCommand struct {
	id bool
	full bool
	chunkSize int64
	silent bool
	fs *flag.FlagSet
}

func (c *GitCleanCommand) FS(fs *flag.FlagSet) {
	c.fs = fs
	fs.BoolVar(&c.id, "id", false, "print only id")
	fs.BoolVar(&c.full, "full", false, "include chunks into manifest")
	fs.Int64Var(&c.chunkSize, "chunk-size", shadow.CHUNK_SIZE,
		"chunk size in bytes")
	fs.BoolVar(&c.silent, "silent", false, "supress warnings")
}

func (c *GitCleanCommand) Do(in io.Reader, out, errOut io.Writer) (err error) {

	s := &shadow.Shadow{}
	if err = s.FromAny(in, c.full, c.chunkSize); err != nil {
		errOut.Write([]byte(err.Error()))
		return
	}
	if s.IsFromShadow && c.silent {
		fmt.Fprintf(errOut, "warning %s is already shadow", c.fs.Args())
	}
	if c.id {
		fmt.Fprintf(out, "%x", s.ID)
	} else {
		err = s.Serialize(out)
	}
	return
}
