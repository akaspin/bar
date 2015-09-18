package cmd
import (
	"flag"
	"io"
	"github.com/akaspin/bar/shadow"
)

type UploadCommand struct {
	endpoint string
	chunkSize int64
	streams int
}

func (c *UploadCommand) FS(fs *flag.FlagSet) {
	fs.StringVar(&c.endpoint, "endpoint", "http://localhost:3000/v1",
		"bard endpoint")
	fs.Int64Var(&c.chunkSize, "chunk-size", shadow.CHUNK_SIZE,
		"upload chunk size")
	fs.IntVar(&c.streams, "streams", 10, "concurrent upload streams count")
}

func (c *UploadCommand) Do(in io.Reader, out, errOut io.Writer) (err error) {

	return
}
