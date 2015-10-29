package fixtures

import (
	"fmt"
	"github.com/nu7hatch/gouuid"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

type Tree struct {
	CWD string
}

func NewTree(where, wd string) *Tree {
	if wd == "" {
		wd, _ = os.Getwd()
	}
	if where == "" {
		id, _ := uuid.NewV4()
		where = id.String()
	}
	root := filepath.Join(wd, "testdata", where)
	os.RemoveAll(root)
	return &Tree{root}
}

func (f *Tree) BlobFilename(n string) string {
	return filepath.Join(f.CWD, n)
}

func (f *Tree) Populate() (err error) {
	if err = os.MkdirAll(f.CWD, 0755); err != nil {
		return
	}
	for _, sub := range []string{"one", "two", "three", ""} {
		for n, s := range map[string]int64{
			"file-one.bin":              3,
			"file-two.bin":              1024*1024*2 + 45,
			"file-three.bin":            1024 + 45,
			"file-four with spaces.bin": 1024 + 45,
		} {
			if err = f.WriteBLOB(filepath.Join(sub, n), s); err != nil {
				return
			}
		}
	}
	//	time.Sleep(time.Millisecond * 100)
	return
}

func (f *Tree) PopulateN(size int64, n int) (err error) {
	for i := 0; i < n; i++ {
		rand.Seed(time.Now().Unix())
		if err = f.WriteBLOB(
			filepath.Join("big", fmt.Sprintf("file-big-%d.bin", i)),
			size+int64(rand.Int31n(100))); err != nil {
			return
		}
	}

	//	time.Sleep(time.Millisecond * 100)
	return
}

func (f *Tree) WriteBLOB(name string, size int64) (err error) {
	err = os.MkdirAll(filepath.Join(f.CWD, filepath.Dir(name)), 0755)
	if err != nil {
		return
	}
	err = MakeNamedBLOB(filepath.Join(f.CWD, name), size)
	return
}

func (f *Tree) KillBLOB(name string) (err error) {
	err = os.Remove(filepath.Join(f.CWD, name))
	return
}

func (f *Tree) Squash() (err error) {
	err = os.RemoveAll(f.CWD)
	return
}
