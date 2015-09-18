package cmd
import (
	"flag"
	"io"
	"github.com/akaspin/bar/shadow"
	"bufio"
	"strings"
	"fmt"
)

// Upload
//
//     find -t file | barctl upload
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
	fs.IntVar(&c.streams, "streams", 10,
		"concurrent upload streams count")
}

func (c *UploadCommand) Do(in io.Reader, out, errOut io.Writer) (err error) {
	// Collect files to upload
	var data []byte
	var feed []struct{
		name string
		id string
	}
	r := bufio.NewReader(in)
	for {
		data, _, err = r.ReadLine()
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			return
		}
		f := strings.Split(string(data), " ")
		ingest := struct{
			name string
			id string
		}{f[0], ""}
		if len(f) > 1 {
			ingest.id = f[1]
		}
		feed = append(feed, ingest)
	}

	// If ids given - do precheck request
	for _, entity := range feed {

		fmt.Println(entity.name, entity.id)
	}

	return
}
