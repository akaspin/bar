package cmd
import (

	"github.com/akaspin/bar/proto/manifest"
	"github.com/akaspin/bar/barc/model"
	"github.com/akaspin/bar/barc/lists"
	"fmt"
	"github.com/akaspin/bar/barc/transport"
)

/*
Replace local shadows with downloaded BLOBs.

	$ barc down my/blobs

 */
type DownCmd struct {
	*BaseSubCommand

	endpoint string
	useGit bool
	maxPool int
	chunkSize int64

	model *model.Model
}

func NewDownCmd(s *BaseSubCommand) SubCommand {
	c := &DownCmd{BaseSubCommand: s}
	c.FS.StringVar(&c.endpoint, "endpoint", "http://localhost:3000/v1",
		"bard endpoint")
	c.FS.BoolVar(&c.useGit, "git", false, "use git infrastructure")
	c.FS.IntVar(&c.maxPool, "pool", 16, "pool size")
	c.FS.Int64Var(&c.chunkSize, "chunk", manifest.CHUNK_SIZE, "chunk size")

	return c
}

func (c *DownCmd) Do() (err error) {
	if c.model, err = model.New(c.WD, c.useGit, c.chunkSize, c.maxPool); err != nil {
		return
	}

	feed := lists.NewFileList(c.FS.Args()...).ListDir(c.WD)

	isDirty, dirty, err := c.model.Check(feed...)
	if err != nil {
		return
	}
	if isDirty {
		err = fmt.Errorf("dirty files in working tree %s", dirty)
		return
	}

	if c.useGit {
		// filter by attrs
		feed, err = c.model.Git.FilterByAttr("bar", feed...)
	}

	blobs, err := c.model.CollectManifests(false, true, feed...)
	if err != nil {
		return
	}

	trans := transport.NewTransport(c.WD, c.endpoint, c.maxPool)
	if err = trans.Download(blobs); err != nil {
		return
	}

	if c.useGit {
		err = c.model.Git.UpdateIndex(blobs.Names()...)
	}

	return
}
