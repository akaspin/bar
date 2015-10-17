package cmd
import (
	"github.com/akaspin/bar/barc/model"
	"fmt"
	"github.com/akaspin/bar/barc/lists"
	"github.com/tamtam-im/logx"
	"github.com/akaspin/bar/barc/transport"
	"flag"
)


/*
This command upload BLOBs to bard and replaces them with shadows.

	$ barctl up my/blobs my/blobs/glob*
*/
type UpCmd struct {
	*Base
	useGit bool
	squash bool

	model *model.Model
}

func NewUpCmd(s *Base) SubCommand {
	c := &UpCmd{Base: s}

	return c
}

func (c *UpCmd) Init(fs *flag.FlagSet) {

	fs.BoolVar(&c.useGit, "git", false, "use git infrastructure")
	fs.BoolVar(&c.squash, "squash", false,
		"replace local BLOBs with shadows after upload")
}

func (c *UpCmd) Do(args []string) (err error) {
	if c.model, err = model.New(c.WD, c.useGit, c.ChunkSize, c.PoolSize); err != nil {
		return
	}

	feed := lists.NewFileList(args...).ListDir(c.WD)

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

	blobs, err := c.model.FeedManifests(true, false, true, feed...)
	if err != nil {
		return
	}

	logx.Debugf("collected blobs %s", blobs.IDMap())

	trans := transport.NewTransport(c.model, "", c.Endpoints, c.PoolSize)

	err = trans.Upload(blobs)
	if err != nil {
		return
	}

	if c.squash {
		if err = c.model.SquashBlobs(blobs); err != nil {
			return
		}
		if c.useGit {
			err = c.model.Git.UpdateIndex(blobs.Names()...)
		}
	}

	return
}

