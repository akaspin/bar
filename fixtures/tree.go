package fixtures
import (
	"os/exec"
	"os"
	"path/filepath"
)


type Tree struct {
	CWD string
}

func NewTree(where string) *Tree {
	if where == "" {
		where, _ = os.Getwd()
		where = filepath.Join(where, "testdata")
	}
	return &Tree{where}
}

func (f *Tree) Populate() (err error) {
	if err = os.MkdirAll(f.CWD, 0755); err != nil {
		return
	}
	for _, sub := range []string{"one", "two", "three", ""} {
		for n, s := range map[string]int64{
			"file-one.bin": 3,
			"file-two.bin": 1024 * 1024 * 2 + 45,
			"file-three.bin": 1024 + 45,
			"file-four with spaces.bin": 1024 + 45,
		} {
			if err = f.WriteBLOB(filepath.Join(sub, n), s); err != nil {
				return
			}
		}
	}
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

func (f *Tree) Run(name string, arg ...string) (res *exec.Cmd) {
	res = exec.Command(name, arg...)
	res.Dir = f.CWD
	return
}

func (f *Tree) InitGit() (err error) {
	_, err = f.Run("git", "init", ).Output()
	return
}



