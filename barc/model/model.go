package model
import (
	"github.com/akaspin/bar/barc/git"
	"io"
	"github.com/akaspin/bar/proto/manifest"
	"time"
	"os"
	"github.com/tamtam-im/logx"
	"sync"
	"fmt"
	"path/filepath"
	"github.com/akaspin/bar/barc/lists"
)


type Model struct {
	WD string
	Git *git.Git
	Hasher *manifest.HasherPool
}

func New(wd string, useGit bool, chunkSize int64, pool int) (res *Model, err error) {
	res = &Model{
		WD: wd,
		Hasher: manifest.NewHasherPool(chunkSize, pool, time.Minute * 30),
	}
	if useGit {
		res.Git, err = git.NewGit(wd)
	}
	return
}

// Check working tree for consistency
func (m *Model) Check(names ...string) (isDirty bool, dirty []string, err error) {
	if m.Git == nil {
		return
	}

	dirty, err = m.Git.DiffFilesWithAttr(names...)
	if err != nil {
		return
	}
	isDirty = len(dirty) > 0
	return
}

// Collect manifests by file names
// Use blobs or/and manifests switches to select specific sources
func (m *Model) CollectManifests(blobs, manifests bool, names ...string) (res lists.Links, err error) {
	res = lists.Links{}
	var errs []error

	resChan := make(chan struct{
		name string
		manifest *manifest.Manifest
	}, 1)
//	}, len(names))
	errChan := make(chan error, 1)

	for _, name := range names {

		go func(n string) {
			h, err1 := m.Hasher.Take()
			if err != nil {
				errChan <- err1
				return
			}
			defer h.Release()

			f, err1 := os.Open(filepath.Join(m.WD, n))
			if err != nil {
				errChan <- err1
				return
			}
			defer f.Close()

			r, isManifest, err1 := h.Peek(f)
			if err != nil {
				errChan <- err1
				return
			}

			if (isManifest && !manifests) || (!isManifest && !blobs) {
				return
			}

			m1, err1 := m.GetManifest(n, r, h)
			if err != nil {
				errChan <- err1
				return
			}
			resChan <- struct{
				name string
				manifest *manifest.Manifest
			}{n, m1}
		}(name)
	}

	for i := 0; i < len(names); i++ {
		select {
		case r1 := <- resChan:
			res[r1.name] = *r1.manifest
		case err1 := <- errChan:
			errs = append(errs, err1)
		}
	}

	if len(errs) > 0 {
		err = fmt.Errorf("errors while collecting manifests %s", errs)
	}
	return
}

func (m *Model) IsBlobs(names []string) (res map[string]bool, err error) {
	res = map[string]bool{}
	wg := sync.WaitGroup{}
	var errs []error

	logx.Tracef("collecting states %s", names)
	for _, n := range names {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			f, err1 := os.Open(filepath.Join(m.WD, name))
			if err1 != nil {
				errs = append(errs, err1)
				return
			}
			defer f.Close()
			_, isManifest, err1 := manifest.Peek(f)
			if err1 != nil {
				errs = append(errs, err1)
				return
			}

			res[name] = !isManifest
		}(n)
	}
	wg.Wait()

	if len(errs) > 0 {
		err = fmt.Errorf("errors %s", errs)
		logx.Error(err)
	}
	return
}

func (m *Model) SquashBlobs(blobs lists.Links) (err error) {
	logx.Tracef("squashing blobs %s", blobs.IDMap())

	wg := sync.WaitGroup{}
	var errs []error
	for n, mt := range blobs {
		wg.Add(1)
		go func(name string, man manifest.Manifest) {
			defer wg.Done()
			absname := filepath.Join(m.WD, name)
			backName := absname + ".bar-backup"
			os.Rename(absname, absname + ".bar-backup")
			os.MkdirAll(filepath.Dir(absname), 0755)
			w, err1 := os.Create(absname)
			if err1 != nil {
				errs = append(errs, err1)
				os.Remove(absname)
				os.Rename(backName, absname)
				return
			}
			err1 = man.Serialize(w)
			if err1 != nil {
				errs = append(errs, err1)
				os.Remove(absname)
				os.Rename(backName, absname)
				return
			}
			defer os.Remove(backName)
			logx.Debugf("squashed %s", name)
		}(n, mt)
	}
	wg.Wait()
	if len(errs) > 0 {
		err = fmt.Errorf("errors while squashing blobs %s", errs)
		return
	}
	logx.Infof("blob %s squashed successfully", blobs.Names())
	return
}


// Get manifest by filename or given reader
func (m *Model) GetManifest(name string, in io.Reader, hasher *manifest.Hasher) (res *manifest.Manifest, err error) {
	if hasher == nil {
		hasher, err = m.Hasher.Take()
		if err != nil {
			return
		}
		defer hasher.Release()

	}

	if in == nil {
		var f *os.File
		if f, err = os.Open(filepath.Join(m.WD, name)); err != nil {
			return
		}
		defer f.Close()
		in = f
	}

	r, isManifest, err := hasher.Peek(in)
	if err != nil {
		return
	}

	if isManifest {
		// ok - just read
		res, err = hasher.MakeFromManifest(r)
		return
	}

	// Hard way. First - try git
	var sideR io.Reader
	if m.Git != nil {
		if sideR = m.getGitReader(name); sideR != nil {
			res, err = hasher.Make(sideR)
			return
		}
	}

	// No git - make from blob
	res, err = hasher.MakeFromBLOB(r)
	return
}

// Try to get reader from git OID
// If git status is dirty or file not in git - just return nil
func (m *Model) getGitReader(name string) (res io.Reader) {
	dirty, _, err := m.Check(name)
	if err != nil {
		logx.Debug(err)
		return
	}
	if dirty {
		err = nil
		logx.Debugf("%s is dirty", name)
		return
	}
	oid, err := m.Git.GetOID(name)
	if err != nil {
		logx.Debug(err)
		return
	}
	res, err = m.Git.Cat(oid)
	if err != nil {
		logx.Debug(err)
		res = nil
	}
	logx.Debugf("manifest for %s parsed from git %s", name, oid)
	return
}

