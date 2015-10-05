package cmd
import (
	"io"
	"github.com/akaspin/bar/proto/manifest"
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
	*BaseSubCommand

	id bool
	chunkSize int64
}

func NewGitCleanCommand(s *BaseSubCommand) SubCommand {
	c := &GitCleanCommand{BaseSubCommand: s}
	c.FS.BoolVar(&c.id, "id", false, "print only id")
	c.FS.Int64Var(&c.chunkSize, "chunk", manifest.CHUNK_SIZE, "preferred chunk size")
	return c
}

func (c *GitCleanCommand) Do() (err error) {
	var name string
	if len(c.FS.Args()) > 0 {
		name = c.FS.Args()[0]
	}

	logx.Debugf("git-clean: %s", name)

	var s *manifest.Manifest

	// check input type
	r, isManifest, err := manifest.Peek(c.Stdin)
	var r2 io.Reader
	if isManifest {
		if s, err = manifest.NewFromManifest(r); err != nil {
			return
		}
		logx.Debugf("%s for %s source is manifest", name)
	} else {
		// blob. Try to check git state
		r2, err = c.getCleanReader(name)
		if err == nil && r2 != nil {
			if s, err = manifest.NewFromManifest(r2); err != nil {
				return
			}
			logx.Debugf("manifest %s for %s created from git cache", s.ID, name)
		} else {
			logx.Debugf("can not get git reader (%s)", err)
			err = nil
			if s, err = manifest.NewFromBLOB(r, c.chunkSize); err != nil {
				return
			}
			logx.Debugf("manifest %s for %s created from BLOB", s.ID, name)
		}
	}

	if c.id {
		fmt.Fprintf(c.Stdout, "%s", s.ID)
	} else {
		err = s.Serialize(c.Stdout)
	}

	return
}

func (c *GitCleanCommand) getCleanReader(name string) (res io.Reader, err error) {
	g, err := git.NewGit(c.WD)
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
	oid, err := g.GetOID(name)
	if err != nil {
		logx.Debugf("error while getting OID for %s", err)
		err = nil
		return
	}
	res, err = g.Cat(oid)
	return
}
