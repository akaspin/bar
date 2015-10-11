package model
import (
	"github.com/akaspin/bar/barc/git"
	"io"
	"github.com/akaspin/bar/proto/manifest"
	"os"
	"github.com/tamtam-im/logx"
	"path/filepath"
	"github.com/akaspin/bar/barc/lists"
	"github.com/akaspin/bar/parmap"
	"time"
)


type Model struct {
	WD string
	Git *git.Git
	Pool *parmap.ParMap
	FdLocks *parmap.LocksPool
	chunkSize int64
}

func New(wd string, useGit bool, chunkSize int64, pool int) (res *Model, err error) {
	res = &Model{
		WD: wd,
		Pool: parmap.NewWorkerPool(pool),
		chunkSize: chunkSize,
		FdLocks: parmap.NewLockPool(64, time.Hour),
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

func (m *Model) ReadChunk(name string, chunk manifest.Chunk, res []byte) (err error) {
	lock, err := m.FdLocks.Take()
	if err != nil {
		return
	}
	defer lock.Release()

	f, err := os.Open(filepath.Join(m.WD, name))
	if err != nil {
		return
	}
	defer f.Close()

	_, err = f.ReadAt(res, chunk.Offset)
	return
}

func (m *Model) FeedManifests(blobs, manifests, strict bool, names ...string) (res lists.Links, err error) {
	req := map[string]interface{}{}
	res = lists.Links{}

	for _, n := range names {
		req[n] = struct{}{}
	}
	res1, err := m.Pool.RunBatch(parmap.Task{
		Map: req,
		Fn: func(name string, trail interface{}) (res interface{}, err error) {
			res, err = m.getManifest(name, blobs, manifests)
			return
		},
		IgnoreErrors: !strict,
	})

	for k, v := range res1 {
		if v.(*manifest.Manifest) != nil {
			res[k] = *(v.(*manifest.Manifest))
		}
	}

	return
}

func (m *Model) getManifest(name string, blobs, manifests bool) (res *manifest.Manifest, err error) {
	lock, err := m.FdLocks.Take()
	if err != nil {
		return
	}
	defer lock.Release()

	f, err := os.Open(filepath.Join(m.WD, name))
	if err != nil {
		return
	}
	defer f.Close()

	var r io.Reader
	var isManifest bool
	if r, isManifest, err = manifest.Peek(f); err != nil {
		return
	}
	if (isManifest && !manifests) || (!isManifest && !blobs) {
		return
	}

	if isManifest {
		res, err = manifest.NewFromManifest(r)
		return
	}
	// Hard way. First - try git
	var sideR io.Reader
	if m.Git != nil {
		if sideR = m.getGitReader(name); sideR != nil {
			res, err = manifest.NewFromAny(sideR, m.chunkSize)
			return
		}
	}
	// No git - make from blob
	res, err = manifest.NewFromBLOB(r, m.chunkSize)
	return
}

func (m *Model) IsBlobs(names ...string) (res map[string]bool, err error) {
	req := map[string]interface{}{}
	res = make(map[string]bool)

	for _, n := range names {
		req[n] = struct{}{}
	}

	logx.Tracef("collecting states %s", names)
	res1, err := m.Pool.RunBatch(parmap.Task{
		Map: req,
		Fn: func(name string, trail interface{}) (res interface{}, err error) {
			lock, err := m.FdLocks.Take()
			if err != nil {
				return
			}
			defer lock.Release()

			f, err := os.Open(filepath.Join(m.WD, name))
			if err != nil {
				return
			}
			defer f.Close()
			_, isManifest, err := manifest.Peek(f)
			if err != nil {
				return
			}
			res = !isManifest
			return
		},
	})
	if err != nil {
		return
	}
	for k, v := range res1 {
		res[k] = v.(bool)
	}

	return
}

func (m *Model) SquashBlobs(blobs lists.Links) (err error) {
	logx.Tracef("squashing blobs %s", blobs.IDMap())

	req := map[string]interface{}{}
	for k, v := range blobs {
		req[k] = v
	}

	if _, err = m.Pool.RunBatch(parmap.Task{
		Map: req,
		Fn: func(name string, in interface{}) (res interface{}, err error) {
			lock, err := m.FdLocks.Take()
			if err != nil {
				return
			}
			defer lock.Release()

			man := in.(manifest.Manifest)
			absname := filepath.Join(m.WD, name)
			backName := absname + ".bar-backup"
			os.Rename(absname, absname + ".bar-backup")
			os.MkdirAll(filepath.Dir(absname), 0755)

			w, err := os.Create(absname)
			if err != nil {
				return
			}
			err = man.Serialize(w)
			if err != nil {
				os.Remove(absname)
				os.Rename(backName, absname)
				return
			}
			defer os.Remove(backName)
			logx.Debugf("squashed %s", name)
			return
		},
		IgnoreErrors: true,
	}); err != nil {
		return
	}

	logx.Infof("blob %s squashed successfully", blobs.Names())
	return
}

// Get manifest by filename or given reader
func (m *Model) GetManifest(name string, in io.Reader) (res *manifest.Manifest, err error) {
	r, isManifest, err := manifest.Peek(in)
	if err != nil {
		return
	}

	if isManifest {
		// ok - just read
		res, err = manifest.NewFromManifest(r)
		return
	}

	// Hard way. First - try git
	var sideR io.Reader
	if m.Git != nil {
		if sideR = m.getGitReader(name); sideR != nil {
			res, err = manifest.NewFromAny(sideR, m.chunkSize)
			return
		}
	}

	// No git - make from blob
	res, err = manifest.NewFromBLOB(r, m.chunkSize)
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

