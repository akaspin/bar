package cmd
import (
	"io"
	"flag"
	"github.com/akaspin/bar/shadow"
	"os"
"io/ioutil"
	"path/filepath"
	"strings"
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
	out, errOut io.Writer
}

func (c *GitSmudgeCmd) Bind(fs *flag.FlagSet, in io.Reader, out, errOut io.Writer) (err error) {
	c.fs = fs
	c.in, c.out, c.errOut = in, out, errOut

	fs.StringVar(&c.endpoint, "endpoint", "http://localhost:3000/v1",
		"bard endpoint")
	fs.BoolVar(&c.shadowChanged, "shadow-changed", false,
		"replace changed blobs with shadows instead download them from bard")
	fs.StringVar(&c.strict, "strict", "",
		"globs to strictly download changed blobs from bard separated by comma")
	return
}

func (c *GitSmudgeCmd) Do() (err error) {
	// read manifest from STDIN
	// We don't need size because where are no BLOBs in
	// staging area.
	ss, err := shadow.New(c.in, 0)
	if err != nil {
		return
	}

	// Check target.
	tmpName, err := c.getTargetShadow(ss)
	if err != nil {
		return
	}

	if tmpName != "" {
		err = c.catFromTmp(tmpName)
		return
	}

	// OK. We has naked shadow in staging area
	if c.isNeedToBeBLOB(c.fs.Args()[0]) {
		// download from bard
	} else {
		err = ss.Serialize(c.out)
	}

	return
}

// Get shadow from target in working tree
// If target is BLOB - also copy it to temporary file
// Returns empty tmpName if target is shadow or absent or
// BLOB what differs from in.
func (c *GitSmudgeCmd) getTargetShadow(in *shadow.Shadow) (tmpName string, err error) {
	n := c.fs.Args()[0]
	info, err := os.Stat(n)
	if err != nil {
		// No file - no problem
		return "", nil
	}

	// Check size to avoid extra ops
	if info.Size() != in.Size {
		return "", nil
	}

	// Size is match
	r, err := os.Open(n)
	if err != nil {
		return
	}
	defer r.Close()

	dr, isShadow, err := shadow.Detect(r)

	if !isShadow {
		// blob - need to copy to temporary file
		var w *os.File
		w, err = ioutil.TempFile("", "bar")
		if err != nil {
			return
		}
		defer w.Close()
		tmpName = w.Name()

		sr := io.TeeReader(dr, w)
		var s *shadow.Shadow
		s, err = shadow.New(sr, info.Size())
		if err != nil {
			defer os.Remove(tmpName)
			return
		}

		if s.ID != in.ID {
			defer os.Remove(tmpName)
			return "", nil
		}
		return tmpName, nil
	} else {
		return "", nil
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
	defer os.Remove(name)

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
