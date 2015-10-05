package cmd
import (
	"github.com/akaspin/bar/proto/manifest"
	"github.com/akaspin/bar/barc/lists"
	"github.com/akaspin/bar/barc/git"
	"sync"
	"os"
	"time"
	"io"
	"fmt"
	"github.com/tamtam-im/logx"
)

/*
Bar spec-out command generates portable specification
and uploads it to bard server. After successful upload
spec-out prints spec URL to STDOUT. All blobs MUST be
uploaded before this command.

	$ barc spec-out
*/
type SpecOutCmd struct {
	*BaseSubCommand

	endpoint string
	useGit bool
	chunkSize int64
	upload bool
	pool int

	git *git.Git
	hasher *manifest.Hasher
}

func NewSpecOutCmd(s *BaseSubCommand) SubCommand  {
	c := &SpecOutCmd{BaseSubCommand: s}
	s.FS.StringVar(&c.endpoint, "endpoint", "http://localhost:3000/v1",
		"bard endpoint")
	s.FS.BoolVar(&c.useGit, "git", false, "use git infrastructure")
	s.FS.Int64Var(&c.chunkSize, "chunk", manifest.CHUNK_SIZE, "preferred chunk size")
	s.FS.BoolVar(&c.upload, "upload", false, "upload spec to bard and print URL")
	s.FS.IntVar(&c.pool, "pool", 16, "pool sizes")
	return c
}

func (c *SpecOutCmd) Do() (err error) {
	if c.useGit {
		if c.git, err = git.NewGit(c.WD); err != nil {
			return
		}
	}
	c.hasher = manifest.NewHasherPool(c.chunkSize, c.pool, time.Minute)

	specmap, err := c.collect()
	logx.Debug(specmap)

	return
}

func (c *SpecOutCmd) collect() (res map[string]manifest.Manifest, err error) {
	names := lists.NewFileList(c.FS.Args()...).ListDir(c.WD)

	logx.Debug(c.WD, names)

	res = map[string]manifest.Manifest{}
	errs := []error{}
	wg := sync.WaitGroup{}
	for _, name := range names {
		wg.Add(1)
		go func(n string) {
			defer wg.Done()
			m, err1 := c.collectOne(n)
			if err1 != nil {
				errs = append(errs, err1)
				return
			}
			res[n] = *m
		}(name)
	}
	wg.Wait()
	if len(errs) > 0 {
		err = fmt.Errorf("errors while collecting manifests %v", errs)
	}
	return
}

// Collect one manifest
func (c *SpecOutCmd) collectOne(name string) (res *manifest.Manifest, err error) {
	f, err := os.Open(name)
	if err != nil {
		return
	}
	defer f.Close()

	r, isManifest, err := manifest.Peek(f)
	if !isManifest {
		// hard way
		if c.useGit {
			r1, err := c.tryFromGit(name)
			if err != nil {
				return nil, err
			}
			if r1 != nil {
				r = r1
			}
			logx.Debugf("got reader from git for %s", name)
		}
	}
	res, err = c.hasher.Make(r)
	return
}

func (c *SpecOutCmd) tryFromGit(name string) (r io.Reader, err error) {
	dirty, err := c.git.DiffFiles(name)
	if err != nil {
		return
	}
	if len(dirty) != 0 {
		err = fmt.Errorf("%s is dirty", name)
		return
	}
	// check by attr
	att, err := c.git.FilterByAttr("bar", name)
	if err != nil {
		return
	}
	if len(att) == 0 {
		return
	}

	oid, err := c.git.GetOID(name)
	if err != nil {
		return
	}
	r, err = c.git.Cat(oid)
	return
}