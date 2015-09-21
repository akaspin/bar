package cmd
import (
	"flag"
	"io"
	"github.com/akaspin/bar/shadow"
	"fmt"
	"github.com/akaspin/bar/barctl/transport"
	"net/url"
	"sync"
	"os"
)

// Upload BLOBS to bard server
//
//     barctl upload FILE [FILE...]
//
// Where FILE is path to regular file.
type UploadCommand struct {
	endpoint string
	chunkSize int64
	streams int
	transportPool *transport.TransportPool
	hasherPool *shadow.HasherPool
	fs *flag.FlagSet
}

func (c *UploadCommand) FS(fs *flag.FlagSet) {
	c.fs = fs
	fs.StringVar(&c.endpoint, "endpoint", "http://localhost:3000/v1",
		"bard endpoint")
	fs.Int64Var(&c.chunkSize, "chunk-size", shadow.CHUNK_SIZE,
		"upload chunk size")
	fs.IntVar(&c.streams, "streams", 10,
		"concurrent upload streams count")
}

func (c *UploadCommand) Do(in io.Reader, out, errOut io.Writer) (err error) {
	u, err := url.Parse(c.endpoint)
	if err != nil {
		return
	}
	c.transportPool = transport.NewTransportPool(u, c.streams, 0)
	c.hasherPool = shadow.NewHasherPool(c.streams, 0, c.chunkSize)

	// Filter files and request existence on bard
	blobShadows := c.collectShadows(errOut)
	fmt.Println(blobShadows)

	return
}

func (c *UploadCommand) collectShadows(errOut io.Writer) (
	res map[string]*shadow.Shadow,
) {
	res = map[string]*shadow.Shadow{}
	var wg sync.WaitGroup
	for _, entity := range c.fs.Args() {
		wg.Add(1)
		go func(entity string) {
			defer wg.Done()
			if err1 := c.collectOneShadow(entity, res); err1 != nil {
				fmt.Fprintln(errOut, err1)
			}
		}(entity)

	}
	wg.Wait()
	return
}

func (c *UploadCommand) collectOneShadow(
	entity string,
	res map[string]*shadow.Shadow,
) (err error) {

	if _, exists := res[entity]; exists {
		return
	}
	var r1 *os.File
	if r1, err = os.Open(entity); err != nil {
		return
	}
	defer r1.Close()

	s, err := c.hasherPool.MakeOne(r1, true)
	if err != nil {
		return
	}
	if !s.IsFromShadow {
		res[entity] = s
	}

	return
}



