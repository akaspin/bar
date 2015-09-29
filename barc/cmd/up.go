package cmd
import (
	"flag"
	"io"
	"os"
	"github.com/akaspin/bar/barc/lists"
	"github.com/akaspin/bar/barc/git"
	"github.com/akaspin/bar/barc/transport"
	"net/url"
	"time"
	"github.com/akaspin/bar/shadow"
	"fmt"
	"sync"
	"github.com/tamtam-im/logx"
)


/*
This command upload BLOBs to bard and replaces
them with shadows.

	$ barctl shadow my/blobs my/blobs/glob*


*/
type UpCmd struct {
	useGit bool
	endpoint string
	poolSize int
	noop bool

	git *git.Git
	tr *transport.TransportPool
	hasher *shadow.HasherPool

	fs *flag.FlagSet
}

func (c *UpCmd) Bind(fs *flag.FlagSet, in io.Reader, out io.Writer) (err error) {
	c.fs = fs

	fs.StringVar(&c.endpoint, "endpoint", "http://localhost:3000/v1",
		"bard endpoint")
	fs.BoolVar(&c.useGit, "git", false, "use git infrastructure")
	fs.BoolVar(&c.noop, "noop", false,
		"do not replace local BLOBs with shadows after upload")
	fs.IntVar(&c.poolSize, "pool", 16, "pool size")
	return
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
	c.hasher = shadow.NewHasherPool(c.poolSize, time.Minute)

	feed, err := c.getFeed()
	if err != nil {
		return
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

	if !c.noop {
		err = c.squash(toSquash)
	}

	return
}

// request bard for existing blobs
func (c *UpCmd) precheck(what map[string]*shadow.Shadow) (res map[string]*shadow.Shadow, err error) {
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

	res = map[string]*shadow.Shadow{}
	for _, name := range idmap {
		res[name] = what[name]
	}
	return
}

func (c *UpCmd) upload(what map[string]*shadow.Shadow) (err error) {
	wg := &sync.WaitGroup{}
	errs := map[string]error{}
	for name, sh := range what {
		wg.Add(1)
		go func(n string, s *shadow.Shadow) {
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

func (c *UpCmd) squash(what map[string]*shadow.Shadow) (err error) {
	wg := sync.WaitGroup{}
	errs := map[string]error{}
	for name, sh := range what {
		wg.Add(1)
		go func(n string, s *shadow.Shadow) {
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

func (c *UpCmd) squashOne(name string, sh *shadow.Shadow) (err error) {
	f, err := os.Create(name)
	if err != nil {
		return
	}
	defer f.Close()
	err = sh.Serialize(f)
	return
}

func (c *UpCmd) collectShadows(in []string) (res map[string]*shadow.Shadow, err error) {
	res = map[string]*shadow.Shadow{}
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
func (c *UpCmd) collectOneShadow(name string) (res *shadow.Shadow, err error) {
	info, err := os.Stat(name)
	if err != nil {
		return
	}

	f, err := os.Open(name)
	if err != nil {
		return
	}
	defer f.Close()

	var r io.Reader
	r, isShadow, err := shadow.Peek(f)
	if isShadow {
		return
	}

	if c.useGit {
		var oid string
		oid, err = c.git.OID(name)
		if err != nil {
			return
		}
		r, err = c.git.Cat(oid)
		if err != nil {
			return
		}
		// using cached manifest - size doesn't matter
	}
	res, err = c.hasher.MakeOne(r, info.Size())

	return
}

func (c *UpCmd) getFeed() (res []string, err error) {
	root, err := os.Getwd()
	if err != nil {
		return
	}

	res, err = lists.NewGlobber(root, c.fs.Args()).List(c.git)
	return
}
