package cmd
import (
	"io"
	"github.com/akaspin/bar/shadow"
	"flag"
	"fmt"
	"os"
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
	fs.BoolVar(&c.silent, "silent", false, "supress warnings")

	return
}

func (c *GitCleanCommand) Do() (err error) {

	info, err := os.Stat(c.fs.Args()[0])
	if err != nil {
		return
	}

	var s *shadow.Shadow
	if s, err = shadow.New(c.in, info.Size()); err != nil {
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
