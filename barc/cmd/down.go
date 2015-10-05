package cmd
import (
	"flag"
	"io"
	"github.com/akaspin/bar/barc/git"
	"github.com/akaspin/bar/barc/transport"
	"github.com/akaspin/bar/shadow"
	"github.com/akaspin/bar/barc/lists"
	"os"
	"sync"
	"fmt"
	"net/url"
	"time"
	"io/ioutil"
)

/*
Replace local shadows with downloaded BLOBs.

	$ barc down my/blobs

 */
type DownCmd struct {
	fs *flag.FlagSet
	wd string

	endpoint string
	useGit bool
	maxPool int
	chunkSize int64

	git *git.Git
	transport *transport.TransportPool
	hasher *shadow.HasherPool

	// Temporary dir
	tmp string
}

func (c *DownCmd) Bind(wd string, fs *flag.FlagSet, in io.Reader, out io.Writer) (err error) {
	c.fs = fs
	c.wd = wd

	fs.StringVar(&c.endpoint, "endpoint", "http://localhost:3000/v1",
		"bard endpoint")
	fs.BoolVar(&c.useGit, "git", false, "use git infrastructure")
	fs.IntVar(&c.maxPool, "pool", 16, "pool size")
	fs.Int64Var(&c.chunkSize, "chunk", shadow.CHUNK_SIZE, "chunk size")
	return
}

func (c *DownCmd) Do() (err error) {
	c.tmp, err = ioutil.TempDir("", "")
	if err != nil {
		return
	}
	defer os.RemoveAll(c.tmp)

	if c.useGit {
		if c.git, err = git.NewGit(c.wd); err != nil {
			return
		}
	}

	u, err := url.Parse(c.endpoint)
	if err != nil {
		return
	}
	c.transport = transport.NewTransportPool(u, c.maxPool, time.Minute)

	// Collect filenames
	feed := lists.NewFileList(c.fs.Args()...).ListDir(c.wd)
	if c.git != nil {
		dirty, err := c.git.DiffFilesWithFilter(feed...)
		if err != nil {
			return err
		}
		if len(dirty) > 0 {
			return fmt.Errorf("dirty files in tree %s", dirty)
		}
	}

	// Collect shadows of needed files
	collected, err := c.collectShadows(feed)

	err = c.download(collected)
	if err != nil {
		return
	}

	if c.useGit {
		err = c.git.UpdateIndex(feed...)
	}
	return
}

// maps BLOB-id to filenames
func (c *DownCmd) download(in map[string]*shadow.Shadow) (err error) {
	res := map[string][]string{}
	errs := []error{}
	for n, s := range in {
		res[s.ID] = append(res[s.ID], n)
	}
	wg := &sync.WaitGroup{}
	for _, names := range res {
		wg.Add(1)
		sh := in[names[0]]
		go func(s *shadow.Shadow, targets []string) {
			defer wg.Done()
			tmp, err1 := c.downloadOne(s)
			if err1 != nil {
				errs = append(errs, err1)
				return
			}
			err1 = c.populate(tmp, targets)
			if err1 != nil {
				errs = append(errs, err1)
				return
			}
		}(sh, names)
	}
	wg.Wait()
	if len(errs) > 0 {
		err = fmt.Errorf("errors while download: %v", errs)
	}
	return
}

// Download one file and return temporary file name
func (c *DownCmd) downloadOne(sh *shadow.Shadow) (res string, err error) {
	f, err := ioutil.TempFile(c.tmp, "")
	if err != nil {
		return
	}
	defer f.Close()

	tr, err := c.transport.Take()
	if err != nil {
		return
	}
	defer c.transport.Release(tr)

	err = tr.GetBLOB(sh.ID, sh.Size, f)
	if err != nil {
		return
	}
	res = f.Name()
	return
}

// Populate temporary file to targets
func (c *DownCmd) populate(t string, targets []string) (err error) {
	wg := sync.WaitGroup{}
	errs := []error{}
	for _, target := range targets {
		wg.Add(1)
		go func(tr string, move bool) {
			defer wg.Done()
			var err1 error
			tName, err1 := c.populateOne(t, target, move)
			if err1 != nil {
				errs = append(errs, err1)
				return
			}
			err1 = os.Rename(tName, target)
			if err1 != nil {
				errs = append(errs, err1)
				return
			}
		}(target, len(targets) == 1)
	}
	wg.Wait()
	if len(errs) > 0 {
		err = fmt.Errorf("errors while populate: %v", errs)
	}
	return
}

func (c *DownCmd) populateOne(src, dst string, move bool) (res string, err error) {
	res = dst + ".bar-tmp"
	if move {
		if err = os.Rename(src, res); err != nil {
			return
		}
	} else {
		var w, r *os.File
		if w, err = os.Create(res); err != nil {
			return
		}
		defer w.Close()
		if r, err = os.Open(src); err != nil {
			return
		}
		defer r.Close()
		if _, err = io.Copy(w, r); err != nil {
			return
		}
	}
	return
}

// Collect shadows
func (c *DownCmd) collectShadows(in []string) (res map[string]*shadow.Shadow, err error) {
	res = map[string]*shadow.Shadow{}
	errs := []error{}

	wg := &sync.WaitGroup{}
	for _, n := range in {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			s, err1 := c.collectOneShadow(name)
			if err1 != nil {
				errs = append(errs, err1)
				return
			}
			if s != nil {
				res[name] = s
			}
		}(n)
	}
	wg.Wait()
	if len(errs) > 0 {
		err = fmt.Errorf("errors while collecting shadows: %v", errs)
	}
	return
}

// Collect one shadow. "res" will be <nil> if file is BLOB.
func (c *DownCmd) collectOneShadow(what string) (res *shadow.Shadow, err error) {
	f, err := os.Open(what)
	if err != nil {
		return
	}
	defer f.Close()

	r, isShadow, err := shadow.Peek(f)
	if err != nil {
		return
	}
	if isShadow {
		// ok it's shadow
		res, err = shadow.NewFromManifest(r)
	}
	return
}
