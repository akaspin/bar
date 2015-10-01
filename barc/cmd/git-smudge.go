package cmd
import (
	"io"
	"flag"
	"github.com/akaspin/bar/shadow"
	"os"
	"github.com/akaspin/bar/barc/transport"
	"net/url"
	"time"
"github.com/tamtam-im/logx"
	"fmt"
)

/*
Git smudge. Used at `git checkout ...`. Equivalent of:

	$ cat staged/shadow-or-BLOB | barctl git-smudge my/file

To use with git register bar in git:

	# .git/config
	[filter "bar"]
		smudge = barctl git-smudge -endpoint=http://my.bar.server/v1 %f

	# .gitattributes
	my/blobs    filter=bar

By default smudge just parse manifest from stdin and pass to stdout

*/
type GitSmudgeCmd struct {
	endpoint string
	smart bool
	chunkSize int64
	maxConn int

	fs *flag.FlagSet
	in io.Reader
	out io.Writer

	timestamp time.Time
}

func (c *GitSmudgeCmd) Bind(wd string, fs *flag.FlagSet, in io.Reader, out io.Writer) (err error) {
	c.fs = fs
	c.in, c.out = in, out
	c.timestamp = time.Now()

	fs.StringVar(&c.endpoint, "endpoint", "http://localhost:3000/v1",
		"bard endpoint")
	fs.BoolVar(&c.smart, "smart", false,
		"download changed blobs from bard")
	fs.Int64Var(&c.chunkSize, "chunk", shadow.CHUNK_SIZE, "preferred chunk size")
	fs.IntVar(&c.maxConn, "pool", 16, "pool size")
	return
}

func (c *GitSmudgeCmd) Do() (err error) {
	name := c.fs.Args()[0]
	logx.Debug(os.Getenv("SMART"))
	if c.smart {
		logx.Debug("smudge in smart mode")
	}

	m, err := shadow.NewFromAny(c.in, c.chunkSize)
	logx.Debugf("smudged manifest for %s (ID: %s from-shadow: %s)",
		name, m.ID, m.IsFromShadow)
	err = m.Serialize(c.out)
	return
}

// download from bard
func (c *GitSmudgeCmd) download(id string, size int64, w io.Writer) (err error) {
	u, err := url.Parse(c.endpoint)
	if err != nil {
		return
	}
	tpool := transport.NewTransportPool(u, c.maxConn, time.Minute)
	t, err := tpool.Take()
	if err != nil {
		return
	}
	defer tpool.Release(t)

	err = t.GetBLOB(id, size, w)
	return
}

// Get shadow from target in working tree
// If target is BLOB - also copy it to temporary file
// Returns empty tmpName if target is shadow or absent or
// BLOB what differs from in.
func (c *GitSmudgeCmd) backupTarget(name string) (res *shadow.Shadow, tmpName string, err error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, "", nil
	}
	defer f.Close()

	dr, isShadow, err := shadow.Peek(f)
	if isShadow {
		// shadow - no problem.
		return nil, "", nil
	}

	tmpName = fmt.Sprintf("%s.bar-%d", name, c.timestamp.UnixNano())

	var w *os.File
	if w, err = os.Create(tmpName); err != nil {
		return
	}
	defer w.Close()
	res, err = c.copyBLOBWithManifest(dr, w)
	return
}

// copy in to out and return manifest
func (c *GitSmudgeCmd) copyBLOBWithManifest(in io.Reader, out io.Writer) (res *shadow.Shadow, err error) {
	inr := io.TeeReader(in, out)
	res, err = shadow.NewFromAny(inr, c.chunkSize)
	return
}
