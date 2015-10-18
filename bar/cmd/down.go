package cmd
import (

	"github.com/akaspin/bar/bar/model"
	"github.com/akaspin/bar/bar/lists"
	"fmt"
	"github.com/akaspin/bar/bar/transport"
	"flag"
)

/*
Replace local shadows with downloaded BLOBs.

	$ bar down my/blobs

 */
type DownCmd struct {
	*Base
	useGit bool
	model *model.Model
}

func NewDownCmd(s *Base) SubCommand {
	c := &DownCmd{Base: s}
	return c
}

func (c *DownCmd) Init(fs *flag.FlagSet) {
	fs.BoolVar(&c.useGit, "git", false, "use git infrastructure")
}

func (c *DownCmd) Do(args []string) (err error) {
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

	blobs, err := c.model.FeedManifests(false, true, true, feed...)
	if err != nil {
		return
	}

	trans := transport.NewTransport(c.model, "", c.Endpoints, c.PoolSize)
	if err = trans.Download(blobs); err != nil {
		return
	}

	if c.useGit {
		err = c.model.Git.UpdateIndex(blobs.Names()...)
	}

	return
}

func (c *DownCmd) Help()  {
	fmt.Fprintln(c.Stderr, "down [OPTIONS] [PATH [PATH] ...]\n")
}

func (c *DownCmd) Summary() {
	fmt.Fprintln(c.Stderr, `
Use bar down to replace manifests in working tree with BLOBs from bard server.
`)
	fmt.Fprintln(c.Stderr, pathTpl)
}

func (c *DownCmd) Description() string {
	return "download blobs from bar server"
}
