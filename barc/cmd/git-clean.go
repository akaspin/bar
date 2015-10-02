package cmd
import (
	"io"
	"github.com/akaspin/bar/shadow"
	"flag"
	"fmt"
	"github.com/tamtam-im/logx"
	"github.com/akaspin/bar/barc/git"
)

/*
Git clean. Used on `git add ...`. Equivalent of:

	$ cat my/file | barctl git-clean my/file

This command takes STDIN and filename and prints shadow
manifest to STDOUT.

To use with git register bar in git:

	# .git/config
	[filter "bar"]
		clean = barc git-clean -endpoint=http://my.bar.server/v1 %f

	# .gitattributes
	my/blobs    filter=bar

STDIN is cat of file in working tree. Git will place STDOUT to stage area.
Optional filename is always relative to git root.
*/
type GitCleanCommand struct {
	id bool
	chunkSize int64

	fs *flag.FlagSet
	in io.Reader
	out io.Writer
	wd string
}

func (c *GitCleanCommand) Bind(wd string, fs *flag.FlagSet, in io.Reader, out io.Writer) (err error) {
	c.fs = fs
	c.in, c.out = in, out
	c.wd = wd

	fs.BoolVar(&c.id, "id", false, "print only id")
	fs.Int64Var(&c.chunkSize, "chunk", shadow.CHUNK_SIZE, "preferred chunk size")
	return
}

func (c *GitCleanCommand) Do() (err error) {
	var name string
	if len(c.fs.Args()) > 0 {
		name = c.fs.Args()[0]
	}

	logx.Debugf("git-clean: %s", name)

	var s *shadow.Shadow

	// check input type
	r, isManifest, err := shadow.Peek(c.in)
	var r2 io.Reader
	if isManifest {
		if s, err = shadow.NewFromManifest(r); err != nil {
			return
		}
		logx.Debugf("%s for %s source is manifest", name)
	} else {
		// blob. Try to check git state
		r2, err = c.getCleanReader(name)
		if err == nil && r2 != nil {
			if s, err = shadow.NewFromManifest(r2); err != nil {
				return
			}
			logx.Debugf("manifest %s for %s created from git cache", s.ID, name)
		} else {
			logx.Debugf("can not get git reader (%s)", err)
			err = nil
			if s, err = shadow.NewFromBLOB(r, c.chunkSize); err != nil {
				return
			}
			logx.Debugf("manifest %s for %s created from BLOB", s.ID, name)
		}
	}

	if c.id {
		fmt.Fprintf(c.out, "%s", s.ID)
	} else {
		err = s.Serialize(c.out)
	}

	return
}

func (c *GitCleanCommand) getCleanReader(name string) (res io.Reader, err error) {
	g, err := git.NewGit(c.wd)
	if err != nil {
		return
	}
	dirty, err := g.DiffFiles(name)
	if err != nil {
		logx.Debugf("error while getting diff %s", err)
		err = nil
		return
	}
	if len(dirty) != 0 {
		// dirty file
		return
	}
	// clean file - read manifest
	oid, err := g.OID(name)
	if err != nil {
		logx.Debugf("error while getting OID for %s", err)
		err = nil
		return
	}
	res, err = g.Cat(oid)
	return
}
