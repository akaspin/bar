package cmd
import (
	"github.com/akaspin/bar/proto/manifest"
	"github.com/akaspin/bar/barc/model"
	"fmt"
	"github.com/akaspin/bar/barc/lists"
	"github.com/tamtam-im/logx"
	"github.com/akaspin/bar/barc/transport"
)


/*
This command upload BLOBs to bard and replaces them with shadows.

	$ barctl up my/blobs my/blobs/glob*
*/
type UpCmd struct {
	*BaseSubCommand

	useGit bool
	endpoint string
	poolSize int
	squash bool
	chunkSize int64

	model *model.Model
}

func NewUpCmd(s *BaseSubCommand) SubCommand {
	c := &UpCmd{BaseSubCommand: s}
	c.FS.StringVar(&c.endpoint, "endpoint", "http://localhost:3000/v1",
		"bard endpoint")
	c.FS.BoolVar(&c.useGit, "git", false, "use git infrastructure")
	c.FS.BoolVar(&c.squash, "squash", false,
		"replace local BLOBs with shadows after upload")
	c.FS.IntVar(&c.poolSize, "pool", 16, "pool size")
	c.FS.Int64Var(&c.chunkSize, "chunk", manifest.CHUNK_SIZE, "preferred chunk size")
	return c
}

func (c *UpCmd) Do() (err error) {
	if c.model, err = model.New(c.WD, c.useGit, c.chunkSize, c.poolSize); err != nil {
		return
	}

	isDirty, dirty, err := c.model.Check()
	if err != nil {
		return
	}
	if isDirty {
		err = fmt.Errorf("dirty files in working tree %s", dirty)
		return
	}

	feed := lists.NewFileList(c.FS.Args()...).ListDir(c.WD)
	if c.useGit {
		// filter by attrs
		feed, err = c.model.Git.FilterByAttr("bar", feed...)
	}

	allBlobs, err := c.model.CollectManifests(true, false, feed...)
	if err != nil {
		return
	}

	logx.Debugf("collected blobs %s", allBlobs.IDMap())

	trans := transport.NewTransport(c.WD, c.endpoint, c.poolSize)

	err = trans.Upload(allBlobs)
	if err != nil {
		return
	}

	if c.squash {
		if err = c.model.SquashBlobs(allBlobs); err != nil {
			return
		}
		err = c.model.Git.UpdateIndex(allBlobs.Names()...)
	}

	return
}

