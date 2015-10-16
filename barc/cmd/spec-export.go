package cmd
import (
	"github.com/akaspin/bar/proto"
	"github.com/akaspin/bar/barc/model"
	"github.com/akaspin/bar/barc/lists"
	"fmt"
	"github.com/akaspin/bar/barc/transport"
	"encoding/json"
	"github.com/tamtam-im/logx"
	"time"
	"os"
	"path/filepath"
	"flag"
)

/*
Export spec to bard
*/
type SpecExportCmd struct {
	*Base

	useGit bool

	upload bool
	doCC bool
	track bool
}

func NewSpecExportCmd(s *Base) SubCommand  {
	c := &SpecExportCmd{Base: s}
	return c
}

func (c *SpecExportCmd) Init(fs *flag.FlagSet) {
	fs.BoolVar(&c.useGit, "git", false, "use git infrastructure")
	fs.BoolVar(&c.upload, "upload", false, "upload spec to bard and print URL")
	fs.BoolVar(&c.doCC, "cc", false, "create spec carbon copy")
}


func (c *SpecExportCmd) Do(args []string) (err error) {
	var mod *model.Model
	if mod, err = model.New(c.WD, c.useGit, c.ChunkSize, c.PoolSize); err != nil {
		return
	}

	feed := lists.NewFileList(c.FS.Args()...).ListDir(c.WD)

	isDirty, dirty, err := mod.Check(feed...)
	if err != nil {
		return
	}
	if isDirty {
		err = fmt.Errorf("dirty files in working tree %s", dirty)
		return
	}

	if c.useGit {
		// filter by attrs
		feed, err = mod.Git.FilterByAttr("bar", feed...)
	}

	blobs, err := mod.FeedManifests(true, true, true, feed...)
	if err != nil {
		return
	}

	// make specmap
	nameMap := map[string]proto.ID{}
	for name, m := range blobs {
		nameMap[name] = m.ID
	}

	spec, err := proto.NewSpec(time.Now().UnixNano(), nameMap, []string{})
	if err != nil {
		return
	}


	if c.doCC {
		ccName := fmt.Sprintf("bar-spec-%d-%s.json",
			time.Now().UnixNano(), spec.ID)
		logx.Infof("storing carbon copy to %s", ccName)
		ccf, err := os.Create(filepath.Join(c.WD, ccName))
		if err != nil {
			return err
		}
		defer ccf.Close()

		if err = json.NewEncoder(ccf).Encode(&spec); err != nil {
			return err
		}
	}

	if !c.upload {
		err = json.NewEncoder(c.Stdout).Encode(&spec)
		return
	}

	trans := transport.NewTransport(mod, "", c.endpoints, c.PoolSize)
	if err = trans.UploadSpec(spec); err != nil {
		return
	}
	fmt.Fprint(c.Stdout, spec.ID)
	return
}