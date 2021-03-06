package model

import (
	"github.com/akaspin/bar/client/git"
	"github.com/akaspin/bar/client/lists"
	"github.com/akaspin/bar/proto"
	"github.com/akaspin/concurrency"
	"github.com/tamtam-im/logx"
	"golang.org/x/net/context"
	"io"
	"os"
	"path/filepath"
	"time"
)

type Model struct {
	WD      string
	Git     *git.Git
	FdLocks *concurrency.Locks
	*concurrency.BatchPool
	chunkSize int64
}

func New(wd string, useGit bool, chunkSize int64, pool int) (res *Model, err error) {
	res = &Model{
		WD:        wd,
		BatchPool: concurrency.NewPool(pool * 32),
		chunkSize: chunkSize,
		FdLocks:   concurrency.NewLocks(context.Background(), pool, time.Minute*5),
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

	delta, err := m.Git.DiffFiles(names...)
	if err != nil {
		return
	}
	if len(delta) == 0 {
		return
	}

	dirty, err = m.Git.FilterByAttr("bar", delta...)
	if err != nil {
		return
	}

	isDirty = len(dirty) > 0
	return
}

func (m *Model) ReadChunk(name string, chunk proto.Chunk, res []byte) (err error) {
	lock, err := m.FdLocks.Take()
	if err != nil {
		return
	}
	defer lock.Release()

	f, err := os.Open(lists.OSFromSlash(lists.OSJoin(m.WD, name)))
	if err != nil {
		return
	}
	defer f.Close()

	_, err = f.ReadAt(res, chunk.Offset)
	return
}

func (m *Model) FeedManifests(blobs, manifests, strict bool, names ...string) (res lists.BlobMap, err error) {
	var req, res1 []interface{}
	for _, n := range names {
		req = append(req, n)
	}

	err = m.BatchPool.Do(
		func(ctx context.Context, in interface{}) (out interface{}, err error) {
			res2, err := m.getManifest(in.(string), blobs, manifests)
			if err != nil {
				return nil, err
			}
			if res2 == nil {
				return
			}
			out = struct {
				name     string
				manifest *proto.Manifest
			}{in.(string), res2}
			return
		},
		&req, &res1, concurrency.DefaultBatchOptions().AllowErrors(),
	)
	res = lists.BlobMap{}
	for _, r := range res1 {
		r1 := r.(struct {
			name     string
			manifest *proto.Manifest
		})
		res[r1.name] = *r1.manifest
	}
	return
}

func (m *Model) getManifest(name string, blobs, manifests bool) (res *proto.Manifest, err error) {
	lock, err := m.FdLocks.Take()
	if err != nil {
		return
	}
	defer lock.Release()

	logx.Tracef("feeding manifest from %s",
		lists.OSFromSlash(lists.OSJoin(m.WD, name)))
	f, err := os.Open(lists.OSFromSlash(lists.OSJoin(m.WD, name)))
	if err != nil {
		logx.Errorf("!!! %s", err)
		return
	}
	defer f.Close()

	var r io.Reader
	var isManifest bool
	if r, isManifest, err = proto.PeekManifest(f); err != nil {
		return
	}
	if (isManifest && !manifests) || (!isManifest && !blobs) {
		return
	}

	if isManifest {
		res, err = proto.NewFromManifest(r)
		return
	}
	// Hard way. First - try git
	var sideR io.Reader
	if m.Git != nil {
		if sideR = m.getGitReader(name); sideR != nil {
			res, err = proto.NewFromAny(sideR, m.chunkSize)
			return
		}
	}
	// No git - make from blob
	res, err = proto.NewFromBLOB(r, m.chunkSize)
	return
}

func (m *Model) IsBlobs(names ...string) (res map[string]bool, err error) {
	res = map[string]bool{}
	var req, res1 []interface{}
	for _, n := range names {
		req = append(req, n)
	}

	err = m.BatchPool.Do(
		func(ctx context.Context, in interface{}) (out interface{}, err error) {
			lock, err := m.FdLocks.Take()
			if err != nil {
				return
			}
			defer lock.Release()

			f, err := os.Open(lists.OSFromSlash(lists.OSJoin(m.WD, in.(string))))
			if err != nil {
				return
			}
			defer f.Close()
			_, isManifest, err := proto.PeekManifest(f)
			if err != nil {
				return
			}
			out = struct {
				name   string
				isBlob bool
			}{in.(string), !isManifest}
			return
		},
		&req, &res1, concurrency.DefaultBatchOptions(),
	)

	for _, r := range res1 {
		r1 := r.(struct {
			name   string
			isBlob bool
		})
		res[r1.name] = r1.isBlob
	}
	return
}

func (m *Model) SquashBlobs(blobs lists.BlobMap) (err error) {
	logx.Tracef("squashing blobs %s", blobs.IDMap())

	var req, res []interface{}
	for _, v := range blobs.ToSlice() {
		req = append(req, v)
	}

	err = m.BatchPool.Do(
		func(ctx context.Context, in interface{}) (out interface{}, err error) {
			r := in.(lists.BlobLink)

			lock, err := m.FdLocks.Take()
			if err != nil {
				return
			}
			defer lock.Release()

			absname := lists.OSFromSlash(lists.OSJoin(m.WD, r.Name))
			backName := absname + ".bar-backup"
			os.Rename(absname, absname+".bar-backup")
			os.MkdirAll(filepath.Dir(absname), 0755)

			w, err := os.Create(absname)
			if err != nil {
				return
			}
			err = r.Manifest.Serialize(w)
			if err != nil {
				os.Remove(absname)
				os.Rename(backName, absname)
				return
			}
			defer os.Remove(backName)
			logx.Debugf("squashed %s", r.Name)
			return
		},
		&req, &res, concurrency.DefaultBatchOptions().AllowErrors(),
	)
	if err != nil {
		return
	}

	logx.Infof("blob %s squashed successfully", blobs.Names())
	return
}

// Get manifest by filename or given reader
func (m *Model) GetManifest(name string, in io.Reader) (res *proto.Manifest, err error) {
	r, isManifest, err := proto.PeekManifest(in)
	if err != nil {
		return
	}

	if isManifest {
		// ok - just read
		res, err = proto.NewFromManifest(r)
		return
	}

	// Hard way. First - try git
	var sideR io.Reader
	if m.Git != nil {
		if sideR = m.getGitReader(name); sideR != nil {
			res, err = proto.NewFromAny(sideR, m.chunkSize)
			return
		}
	}

	// No git - make from blob
	res, err = proto.NewFromBLOB(r, m.chunkSize)
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

func (m *Model) Close() {
	m.FdLocks.Close()
	m.BatchPool.Close()
}
