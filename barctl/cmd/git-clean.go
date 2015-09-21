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
	in io.Reader
	out, errOut io.Writer
}

func (c *GitCleanCommand) Bind(fs *flag.FlagSet, in io.Reader, out, errOut io.Writer) (err error) {
	c.fs = fs
	c.in, c.out, c.errOut = in, out, errOut

	fs.BoolVar(&c.id, "id", false, "print only id")
	fs.BoolVar(&c.full, "full", false, "include chunks into manifest")
	fs.Int64Var(&c.chunkSize, "chunk-size", shadow.CHUNK_SIZE,
		"chunk size in bytes")
	fs.BoolVar(&c.silent, "silent", false, "supress warnings")

	return
}

func (c *GitCleanCommand) Do() (err error) {

	s := &shadow.Shadow{}
	if err = s.FromAny(c.in, c.full, c.chunkSize); err != nil {
		c.errOut.Write([]byte(err.Error()))
		return
	}
	if s.IsFromShadow && c.silent {
		fmt.Fprintf(c.errOut, "warning %s is already shadow", c.fs.Args())
	}
	if c.id {
		fmt.Fprintf(c.out, "%s", s.ID)
	} else {
		err = s.Serialize(c.out)
	}
	return
}
