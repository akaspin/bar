package cmd
import (
	"io"
	"flag"
	"github.com/akaspin/bar/shadow"
	"os"
	"path/filepath"
	"strings"
	"github.com/akaspin/bar/barc/transport"
	"net/url"
	"time"
"github.com/tamtam-im/logx"
	"fmt"
)

/*
Git smudge. Used at `git checkout ...`. Equivalent of:

	$ cat staged/shadow | barctl git-smudge my/file

This command takes cat of shadow manifest as STDIN and FILE and
writes file contents to STDOUT.

To use with git register bar in git:

	# .git/config
	[filter "bar"]
		smudge = barctl git-smudge -endpoint=http://my.bar.server/v1 %f

	# .gitattributes
	my/blobs    filter=bar

STDIN is cat of file in git stage area. Git will place STDOUT to
working tree. If FILE is BLOB. git-smudge will check it ID and download
correct revision from bard. To disable this behaviour and replace changed
BLOBs with shadows use -shadow-changed option.
*/
type GitSmudgeCmd struct {
	endpoint string
	shadowChanged bool
	strict string

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
	fs.StringVar(&c.strict, "strict", "",
		"globs to strictly download changed blobs from bard separated by comma")
	return
}

func (c *GitSmudgeCmd) Do() (err error) {

	// Check target.
	wtid, tmpName, err := c.getTargetShadow()
	if err != nil {
		return
	}

	// read manifest from STDIN
	// We don't need size because where are no BLOBs in
	// staging area.
	ss, err := shadow.New(c.in, 0)
	if err != nil {
		if tmpName != "" {
			logx.Errorf(
				"can not parse manifest from staging area: %s. old blob stored in %s",
				err, c.tmpName())
		} else {
			logx.Errorf("can not parse manifest from staging area: %s", err)
		}
		return
	}

	if tmpName != "" && wtid == ss.ID {
		logx.Debugf("BLOB %s not changed")
		err = c.catFromTmp(tmpName)
		defer os.Remove(tmpName)
		return
	}

	fname := c.fs.Args()[0]

	// OK. We has naked shadow in staging area
	if c.isNeedToBeBLOB(fname) {
		// download from bard
		logx.Infof("downloading %s (%d bytes, %s)", c.fs.Args()[0], ss.Size, ss.ID)

		u, err := url.Parse(c.endpoint)
		if err != nil {
			logx.Error(err)
			return err
		}
		trPool := transport.NewTransportPool(u, 10, time.Minute * 5)
		if err != nil {
			logx.Error(err)
			return err
		}
		tr, err := trPool.Take()
		if err != nil {
			logx.Error(err)
			return err
		}
		defer trPool.Release(tr)

		err = tr.GetBLOB(ss.ID, ss.Size, c.out)
		if err != nil {
			logx.Error(
				"can not download %s (%s) from bard: %s. old blob stored in %s",
				fname, ss.ID, err, fname + ".bar-" + ss.ID)
			return err
		} else {
			defer os.Remove(fname + ".bar-" + ss.ID)
		}
	} else {
		err = ss.Serialize(c.out)
	}

	return
}

// Get shadow from target in working tree
// If target is BLOB - also copy it to temporary file
// Returns empty tmpName if target is shadow or absent or
// BLOB what differs from in.
func (c *GitSmudgeCmd) getTargetShadow() (id string, tmpName string, err error) {
	n := c.fs.Args()[0]
	info, err := os.Stat(n)
	if err != nil {
		// No file - no problem
		return "", "", nil
	}

	// Size is match
	r, err := os.Open(n)
	if err != nil {
		return
	}
	defer r.Close()

	dr, isShadow, err := shadow.Peek(r)

	if !isShadow {
		// blob - need to copy to temporary file
		var w *os.File
		w, err = os.Create(c.tmpName())
		if err != nil {
			return
		}
		defer w.Close()
		tmpName = w.Name()

		sr := io.TeeReader(dr, w)
		var s *shadow.Shadow
		s, err = shadow.New(sr, info.Size())
		if err != nil {
			logx.Error(err)
			return
		}

		return s.ID, tmpName, nil
	} else {
		return "", "", nil
	}

	return
}

// Cat file to given writer and remove it
func (c *GitSmudgeCmd) catFromTmp(name string) (err error) {
	r, err := os.Open(name)
	if err != nil {
		return
	}
	defer r.Close()

	_, err = io.Copy(c.out, r)
	return
}

// Check target filename by globs in strict
func (c *GitSmudgeCmd) isNeedToBeBLOB(name string) (res bool) {
	for _, g := range strings.Split(c.strict, ",") {
		if res, _ = filepath.Match(g, name); res {
			return
		}
	}
	return
}

func (c *GitSmudgeCmd) tmpName() string {
	return fmt.Sprintf("%s.bar-%d", c.fs.Args()[0], c.timestamp.UnixNano())
}
