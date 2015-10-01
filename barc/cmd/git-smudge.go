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

*/
type GitSmudgeCmd struct {
	endpoint string
	shadowChanged bool
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
	fs.BoolVar(&c.shadowChanged, "shadow-changed", false,
		"replace changed blobs with shadows instead download them from bard")
	fs.Int64Var(&c.chunkSize, "chunk", shadow.CHUNK_SIZE, "preferred chunk size")
	fs.IntVar(&c.maxConn, "pool", 16, "pool size")
	return
}

func (c *GitSmudgeCmd) Do() (err error) {
	fname := c.fs.Args()[0]
	var backupShadow, srcShadow *shadow.Shadow

	// Check target.
	backupShadow, backupName, err := c.backupTarget(fname)
	if err != nil {
		return
	}

	in, isManifest, err := shadow.Peek(c.in)
	if !isManifest {
		srcShadow, err = c.copyBLOBWithManifest(in, c.out)
		logx.Warningf(
			"%s (id %s) is BLOB in staging area and maybe not exists on %s",
			fname, srcShadow.ID, c.endpoint,
		)
	} else {
		// ok source is manifest
		if srcShadow, err = shadow.NewFromManifest(in); err != nil {
			return
		}
		if backupShadow != nil {
			// target is blob
			if backupShadow.ID == srcShadow.ID {
				// just cat tempfile
				logx.Debugf("source id %s is ok for target BLOB %s",
					backupShadow.ID, fname)

				var f *os.File
				f, err = os.Open(backupName)
				if err != nil {
					return err
				}
				defer f.Close()
				if _, err = io.Copy(c.out, f); err != nil {
					return
				}
				defer os.Remove(backupName)
			} else {
				logx.Infof("source manifest %s differs target blob %s %s",
					srcShadow.ID, fname, backupShadow.ID)
				// target differs src
				if c.shadowChanged {
					// just cat in to out
					logx.Warningf(
						"shadow-changed enabled. storing %s as manifest %s",
						fname, srcShadow.ID)
					if err = srcShadow.Serialize(c.out); err != nil {
						return
					}
					defer os.Remove(backupName)
				} else {
					// need to download from bard
					logx.Infof("downloading %s for %s", srcShadow.ID, fname)
					if err = c.download(srcShadow.ID, srcShadow.Size, c.out); err != nil {
						return
					}
					defer os.Remove(backupName)
				}
			}
		} else {
			// target is shadow - just cat input
			logx.Debugf("target is manifest")
			if srcShadow, err = shadow.NewFromManifest(in); err != nil {
				return
			}
			if err = srcShadow.Serialize(c.out); err != nil {
				return
			}
		}

	}
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
