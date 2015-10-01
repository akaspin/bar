package cmd
import (
	"flag"
	"github.com/akaspin/bar/shadow"
	"io"
	"fmt"
	"os"
	"github.com/akaspin/bar/barc/git"
	"github.com/tamtam-im/logx"
)

/*
Cat for commit diff.

	$ git diff --cached --staged
    ...
	+BAR-SHADOW-BLOB 859a7a7603028deeb3b66234cffa5191466d1a0538e449a19812273b0d98dc1c

This command used as textconv and always invoked from root of git repo
*/
type GitTextconvCmd struct {
	fs *flag.FlagSet
	out io.Writer

	strict bool
}

func (c *GitTextconvCmd) Bind(wd string, fs *flag.FlagSet, in io.Reader, out io.Writer) (err error) {
	c.fs = fs
	c.out = out
	return
}

func (c *GitTextconvCmd) Do() (err error) {
	n := c.fs.Args()[0]
	logx.Debugf("cat %s", n)

	g, err := git.NewGit("")
	if err != nil {
		return
	}

	fr, err := os.Open(n)
	if err != nil {
		return
	}
	defer fr.Close()

	r, isShadow, err := shadow.Peek(fr)
	if err != nil {
		err = nil
		return
	}

	var s *shadow.Shadow

	if !isShadow {
		// try to get manifest from git index
		logx.Debugf("%s is BLOB trying to take manifest from git index", n)
		var oid string
		oid, err = g.OID(n)
		if err != nil {
			logx.Debugf("can not find %s in git index", n)
//			return
		} else {
			logx.Debugf("taking manifest from %s", oid)
			r, err = g.Cat(oid)
			if err != nil {
				return
			}

		}
	}
	s, err = shadow.NewFromAny(r, shadow.CHUNK_SIZE)

	if err != nil {
		return
	}
	fmt.Fprintf(c.out, "BAR-SHADOW-BLOB %s\n", s.ID)
	return
}

