package model
import (
	"github.com/akaspin/bar/barc/git"
	"io"
	"github.com/akaspin/bar/proto/manifest"
	"time"
	"os"
	"github.com/tamtam-im/logx"
)


type Model struct {
	WD string
	Git *git.Git
	Hasher *manifest.Hasher
}

func New(wd string, useGit bool, chunkSize int64, pool int) (res *Model, err error) {
	res = &Model{
		WD: wd,
		Hasher: manifest.NewHasherPool(chunkSize, pool, time.Minute * 5),
	}
	if useGit {
		res.Git, err = git.NewGit(wd)
	}
	return
}

// Check working tree for consistency
func (m *Model) Check(names ...string) (dirty bool, err error) {
	if m.Git == nil {
		return
	}

	d, err := m.Git.DiffFilesWithAttr(names...)
	if err != nil {
		return
	}
	dirty = len(d) > 0
	return
}

// Get manifest by filename or given reader
func (m *Model) GetManifest(name string, in io.Reader) (res *manifest.Manifest, err error) {
	if in == nil {
		var f *os.File
		if f, err = os.Open(name); err != nil {
			return
		}
		defer f.Close()
		in = f
	}

	r, isManifest, err := manifest.Peek(in)
	if err != nil {
		return
	}

	if isManifest {
		// ok - just read
		res, err = m.Hasher.Make(r)
		return
	}

	// Hard way. First - try git
	var sideR io.Reader
	if m.Git != nil {
		if sideR = m.getGitReader(name); sideR != nil {
			res, err = m.Hasher.Make(sideR)
			return
		}
	}

	// No git - make from blob
	res, err = m.Hasher.Make(r)
	return
}

// Try to get reader from git OID
// If git status is dirty or file not in git - just return nil
func (m *Model) getGitReader(name string) (res io.Reader) {
	dirty, err := m.Check(name)
	if err != nil {
		logx.Debug(err)
		return
	}
	if dirty {
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
	return
}