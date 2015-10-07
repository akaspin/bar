package cmd
import (
	"github.com/akaspin/bar/proto/manifest"
	"github.com/akaspin/bar/barc/model"
	"github.com/akaspin/bar/barc/lists"
	"fmt"
	"github.com/akaspin/bar/proto"
	"github.com/akaspin/bar/barc/transport"
	"encoding/json"
)

/*
Export spec to bard
*/
type SpecExportCmd struct {
	*BaseSubCommand

	endpoint string
	useGit bool
	chunkSize int64
	upload bool
	pool int

}

func NewSpecExportCmd(s *BaseSubCommand) SubCommand  {
	c := &SpecExportCmd{BaseSubCommand: s}
	s.FS.StringVar(&c.endpoint, "endpoint", "http://localhost:3000/v1",
		"bard endpoint")
	s.FS.BoolVar(&c.useGit, "git", false, "use git infrastructure")
	s.FS.Int64Var(&c.chunkSize, "chunk", manifest.CHUNK_SIZE, "preferred chunk size")
	s.FS.BoolVar(&c.upload, "upload", false, "upload spec to bard and print URL")
	s.FS.IntVar(&c.pool, "pool", 16, "pool sizes")
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

	blobs, err := mod.CollectManifests(true, true, feed...)
	if err != nil {
		return
	}

	// make specmap
	nameMap := map[string]string{}
	for name, m := range blobs {
		nameMap[name] = m.ID
	}

	spec, err := proto.NewSpec(nameMap)
	if err != nil {
		return
	}

	if !c.upload {
		err = json.NewEncoder(c.Stdout).Encode(&spec)
		return
	}

	trans := transport.NewTransport(c.WD, c.endpoint, c.pool)
	if err = trans.UploadSpec(spec); err != nil {
		return
	}
	fmt.Fprintf(c.Stdout, spec.ID)
	return
}