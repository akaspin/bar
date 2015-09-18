package cmd
import (
	"io"
	"github.com/akaspin/bar/shadow"
	"flag"
)

type CleanCommand struct {
	full bool
	chunkSize int64
}

func (c *CleanCommand) FS(fs *flag.FlagSet) {
	fs.BoolVar(&c.full, "full", false, "include chunks into manifest")
	fs.Int64Var(&c.chunkSize, "chunk-size", shadow.CHUNK_SIZE,
		"chunk size in bytes")
}

func (c *CleanCommand) Do(in io.Reader, out, errOut io.Writer) (err error) {

	s := &shadow.Shadow{}
	if err = s.FromAny(in, c.full, c.chunkSize); err != nil {
		errOut.Write([]byte(err.Error()))
		return
	}
	err = s.Serialize(out)
	return
}
