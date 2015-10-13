package cmd
import (
	"github.com/akaspin/bar/proto/manifest"
	"github.com/akaspin/bar/barc/model"
	"github.com/akaspin/bar/barc/lists"
	"fmt"
	"github.com/akaspin/bar/proto"
	"github.com/akaspin/bar/barc/transport"
	"encoding/json"
	"github.com/tamtam-im/logx"
	"time"
	"os"
	"path/filepath"
)

/*
Export spec to bard
*/
type SpecExportCmd struct {
	*BaseSubCommand

	httpEndpoint string
	rpcEndpoints string

	useGit bool
	chunkSize int64
	pool int

	upload bool
	doCC bool
	track bool
}

func NewSpecExportCmd(s *BaseSubCommand) SubCommand  {
	c := &SpecExportCmd{BaseSubCommand: s}
	c.FS.StringVar(&c.httpEndpoint, "http", "http://localhost:3000/v1",
		"bard http endpoint")
	c.FS.StringVar(&c.rpcEndpoints, "rpc", "localhost:3001",
		"bard rpc endpoints separated by comma")
	s.FS.BoolVar(&c.useGit, "git", false, "use git infrastructure")
	s.FS.Int64Var(&c.chunkSize, "chunk", manifest.CHUNK_SIZE, "preferred chunk size")
	s.FS.IntVar(&c.pool, "pool", 16, "pool sizes")

	s.FS.BoolVar(&c.upload, "upload", false, "upload spec to bard and print URL")
	s.FS.BoolVar(&c.doCC, "cc", false, "create spec carbon copy")
	return c
}

func (c *SpecExportCmd) Do() (err error) {
	var mod *model.Model
	if mod, err = model.New(c.WD, c.useGit, c.chunkSize, c.pool); err != nil {
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
	nameMap := map[string]string{}
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

	trans := transport.NewTransport(mod, c.httpEndpoint, c.rpcEndpoints, c.pool)
	if err = trans.UploadSpec(spec); err != nil {
		return
	}
	fmt.Fprintf(c.Stdout, spec.ID)
	return
}