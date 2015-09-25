package cmd
import (
	"flag"
	"io"
	"github.com/akaspin/bar/shadow"
	"github.com/akaspin/bar/barctl/transport"
	"net/url"
	"sync"
	"os"
"github.com/tamtam-im/logx"
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
	in io.Reader
	out io.Writer
}

func (c *UploadCommand) Bind(fs *flag.FlagSet, in io.Reader, out io.Writer) (err error) {
	c.fs = fs
	c.in = in
	c.out = out

	fs.StringVar(&c.endpoint, "endpoint", "http://localhost:3000/v1",
		"bard endpoint")
	fs.Int64Var(&c.chunkSize, "chunk-size", shadow.CHUNK_SIZE,
		"upload chunk size")
	fs.IntVar(&c.streams, "streams", 10,
		"concurrent upload streams count")

	return
}

func (c *UploadCommand) Do() (err error) {
	u, err := url.Parse(c.endpoint)
	if err != nil {
		return
	}
	c.transportPool = transport.NewTransportPool(u, c.streams, 0)
	c.hasherPool = shadow.NewHasherPool(c.streams, 0)

	toUpload, err := c.precheck()

	var wg sync.WaitGroup
	for f, s := range toUpload {
		wg.Add(1)
		go func(f string, s *shadow.Shadow) {
			defer wg.Done()
			t, err := c.transportPool.Take()
			if err != nil {
				logx.Error(err)
				return
			}
			defer c.transportPool.Release(t)
			err = t.Push(f, s)
			if err != nil {
				logx.Error(err)
			}
		}(f, s)
	}
	wg.Wait()

	return
}

func (c *UploadCommand) precheck() (res map[string]*shadow.Shadow, err error) {
	// Filter files and request existence on bard
	res = c.collectShadows()

	// Precheck on bard
	var req []string
	t, err := c.transportPool.Take()
	if err != nil {
		return
	}
	defer c.transportPool.Release(t)
	for _, s := range res {
		req = append(req, s.ID)
	}
	resp, err := t.Check(req)
	if err != nil {
		return
	}

	byID := map[string]string{}
	for f, s := range res {
		if _, exists := byID[s.ID]; !exists {
			byID[s.ID] = f
		}
	}
	for _, id := range resp {
		if f, ok := byID[id]; ok {
			delete(res, f)
		}
	}

	return

}

func (c *UploadCommand) collectShadows() (
	res map[string]*shadow.Shadow,
) {
	res = map[string]*shadow.Shadow{}
	var wg sync.WaitGroup
	for _, entity := range c.fs.Args() {
		wg.Add(1)
		go func(entity string) {
			defer wg.Done()
			if err1 := c.collectOneShadow(entity, res); err1 != nil {
				logx.Error(err1)
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

	s, err := c.hasherPool.MakeOne(r1, 0)
	if err != nil {
		return
	}
	if !s.IsFromShadow {
		res[entity] = s
	}

	return
}



