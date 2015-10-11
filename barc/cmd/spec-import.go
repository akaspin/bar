package cmd
import (
	"github.com/akaspin/bar/proto/manifest"
	"github.com/akaspin/bar/barc/lists"
	"encoding/json"
	"github.com/akaspin/bar/barc/transport"
	"github.com/akaspin/bar/barc/model"
	"github.com/tamtam-im/logx"
	"fmt"
	"path/filepath"
)


/*
Import spec from bard and populate manifests

	$ cat spec.json | bar spec-import -raw
	$ bar spec-import 1bcaa5...578bd24
	$ bar spec-import http://localhost:3000/v1/spec/1bcaa5...578bd24
	$ bar spec-import http://localhost:3000/v1/spec/1bcaa5...578bd24.json

*/
type SpecImportCmd struct  {
	*BaseSubCommand

	endpoint string
	useGit bool
	chunkSize int64
	pool int

	raw bool
	fetchBlobs bool
}

func NewSpecImportCmd(s *BaseSubCommand) SubCommand  {
	c := &SpecImportCmd{BaseSubCommand: s}
	s.FS.StringVar(&c.endpoint, "endpoint", "http://localhost:3000/v1",
		"bard endpoint")
	s.FS.BoolVar(&c.useGit, "git", false, "use git infrastructure")
	s.FS.Int64Var(&c.chunkSize, "chunk", manifest.CHUNK_SIZE, "preferred chunk size")
	s.FS.IntVar(&c.pool, "pool", 16, "pool sizes")
	s.FS.BoolVar(&c.raw, "raw", false, "read spec from STDIN")
	s.FS.BoolVar(&c.fetchBlobs, "blobs", false, "download blobs")
	return c
}

func (c *SpecImportCmd) Do() (err error) {
	// get spec first
	var spec lists.Links
	mod, err := model.New(c.WD, c.useGit, c.chunkSize, c.pool)
	if err != nil {
		return
	}
	trans := transport.NewTransport(mod, c.endpoint, c.pool)

	if c.raw {
		if err = json.NewDecoder(c.Stdin).Decode(&spec); err != nil {
			return
		}
	} else {
		// tree spec types
		id := c.FS.Arg(0)

		if spec, err = trans.GetSpec(id); err != nil {
			logx.Debug(spec, err)
			return
		}
	}

	names := spec.Names()

	logx.Debugf("importing %s", names)

	if c.useGit {
		// If git is used - check names for attrs
		byAttr, err := mod.Git.FilterByAttr("bar", names...)
		if err != nil {
			return err
		}

		diff := []string{}
		attrs := map[string]struct{}{}
		for _, x := range byAttr {
			attrs[x] = struct{}{}
		}

		for _, x := range names {
			if _, ok := attrs[x]; !ok {
				diff = append(diff, x)
			}
		}
		if len(diff) > 0 {
			return fmt.Errorf("some spec blobs is not under bar control %s", diff)
		}
	}

	// get stored links, ignore errors
	stored, _ := mod.FeedManifests(true, true, false, names...)

	logx.Debugf("already stored %s", stored.Names())

	// squash present
	toSquash := lists.Links{}
	for n, m := range spec {
		m1, ok := stored[filepath.FromSlash(n)]
		if !ok || m.ID != m1.ID {
			toSquash[n] = spec[n]
		}
	}

	if err = mod.SquashBlobs(toSquash); err != nil {
		return
	}

	for k, _ := range spec {
		fmt.Fprintf(c.Stdout, "%s ", filepath.FromSlash(k))
	}
	return
}

