package cmd
import (
	"io"
	"os"
	"github.com/akaspin/bar/barc/lists"
	"github.com/akaspin/bar/barc/git"
	"github.com/akaspin/bar/barc/transport"
	"net/url"
	"time"
	"github.com/akaspin/bar/proto/manifest"
	"fmt"
	"sync"
	"github.com/tamtam-im/logx"
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

	git *git.Git
	tr *transport.TransportPool
	hasher *manifest.Hasher
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
	if c.useGit {
		if c.git, err = git.NewGit(""); err != nil {
			return
		}
	}

	u, err := url.Parse(c.endpoint)
	if err != nil {
		return
	}

	c.tr = transport.NewTransportPool(u, c.poolSize, time.Minute)
	c.hasher = manifest.NewHasherPool(c.chunkSize, c.poolSize, time.Minute)

	feed := lists.NewFileList(c.FS.Args()...).ListDir(c.WD)
	if c.git != nil {
		dirty, err := c.git.DiffFilesWithAttr(feed...)
		if err != nil {
			return err
		}
		if len(dirty) > 0 {
			return fmt.Errorf("dirty files in tree %s", dirty)
		}
	}

	logx.Debugf("files to upload %s", feed)

	toSquash, err := c.collectShadows(feed)
	if err != nil {
		return
	}

	toUpload, err := c.precheck(toSquash)
	if err != nil {
		return
	}

	err = c.upload(toUpload)
	if err != nil {
		return
	}

	if c.squash {
		err = c.squashBLOBs(toSquash)
	}

	return
}

// request bard for existing blobs
func (c *UpCmd) precheck(what map[string]*manifest.Manifest) (res map[string]*manifest.Manifest, err error) {
	idmap := map[string]string{}
	req := []string{}
	for name, sh := range what {
		idmap[sh.ID] = name
		req = append(req, sh.ID)
	}
	tr, err := c.tr.Take()
	if err != nil {
		return
	}
	defer c.tr.Release(tr)

	resp, err := tr.Check(req)
	if err != nil {
		return
	}

	for _, id := range resp {
		delete(idmap, id)
	}

	res = map[string]*manifest.Manifest{}
	for _, name := range idmap {
		res[name] = what[name]
	}
	return
}

func (c *UpCmd) upload(what map[string]*manifest.Manifest) (err error) {
	wg := &sync.WaitGroup{}
	errs := map[string]error{}
	for name, sh := range what {
		wg.Add(1)
		go func(n string, s *manifest.Manifest) {
			defer wg.Done()
			tr, err1 := c.tr.Take()
			if err1 != nil {
				errs[n] = err1
				return
			}
			defer c.tr.Release(tr)

			err1 = tr.Push(n, s)
			if err1 != nil {
				errs[n] = err1
			}
		}(name, sh)
	}
	wg.Wait()
	if len(errs) > 0 {
		err = fmt.Errorf("errors while upload: %v", errs)
	}
	return
}

func (c *UpCmd) squashBLOBs(what map[string]*manifest.Manifest) (err error) {
	wg := sync.WaitGroup{}
	errs := map[string]error{}
	for name, sh := range what {
		wg.Add(1)
		go func(n string, s *manifest.Manifest) {
			defer wg.Done()
			err1 := c.squashOne(n, s)
			if err1 != nil {
				errs[n] = err1
			}
		}(name, sh)
	}
	wg.Wait()
	if len(errs) > 0 {
		err = fmt.Errorf("errors while squash: %v", errs)
		return
	}
	// update git index
	if c.useGit {
		var toReadd []string
		for n, _ := range what {
			toReadd = append(toReadd, n)
		}
		// TODO: use git update-index
		err = c.git.UpdateIndex(toReadd...)
	}
	return
}

func (c *UpCmd) squashOne(name string, sh *manifest.Manifest) (err error) {
	f, err := os.Create(name)
	if err != nil {
		return
	}
	defer f.Close()
	err = sh.Serialize(f)
	return
}

func (c *UpCmd) collectShadows(in []string) (res map[string]*manifest.Manifest, err error) {
	res = map[string]*manifest.Manifest{}
	errs := map[string]error{}
	wg := &sync.WaitGroup{}
	for _, n := range in {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			sh, err1 := c.collectOneShadow(name)
			if err1 != nil {
				errs[name] = err1
			} else if sh != nil {
				res[name] = sh
			}
		}(n)
	}
	wg.Wait()
	if len(errs) > 0 {
		err = fmt.Errorf("errors while collect shadows: %v", errs)
	}
	return
}

// Collect shadow by filename
// Returns nil if file is already shadow
func (c *UpCmd) collectOneShadow(name string) (res *manifest.Manifest, err error) {
	f, err := os.Open(name)
	if err != nil {
		return
	}
	defer f.Close()

	var r io.Reader
	r, isShadow, err := manifest.Peek(f)
	if isShadow {
		return
	}

	if c.useGit {
		var oid string
		oid, err = c.git.GetOID(name)
		if err != nil {
			return
		}
		r, err = c.git.Cat(oid)
		if err != nil {
			return
		}
		// using cached manifest - size doesn't matter
	}
	res, err = c.hasher.Make(r)

	return
}

