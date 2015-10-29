package lists

import (
	"path/filepath"
)

/*
Mapper remaps paths from primary to secondary or vise-versa.

	/home
		/dir        < root
			/one    < cwd
			/two

	my/path -> to root: one/my/path
	one/my/path -> to cwd: my/path

Useful to remap paths to root directory
*/
type Mapper struct {

	// Primary root
	CWD string

	// Secondary root
	Root string
}

func NewMapper(cwd, root string) (res *Mapper) {
	res = &Mapper{cwd, root}
	return
}

// Takes paths relative to cwd and returns paths relative
// to root.
func (m *Mapper) ToRoot(arg ...string) (res []string, err error) {
	res, err = m.remap(m.CWD, m.Root, arg...)
	return
}

func (m *Mapper) FromRoot(arg ...string) (res []string, err error) {
	res, err = m.remap(m.Root, m.CWD, arg...)
	return
}

func (m *Mapper) remap(from, to string, arg ...string) (res []string, err error) {
	if from == to {
		res = arg
		return
	}
	var one string
	for _, p := range arg {
		if one, err = filepath.Rel(to, filepath.Join(from, p)); err != nil {
			return
		}
		res = append(res, filepath.Clean(one))
	}
	return
}

func (m *Mapper) ToShell(arg ...string) (res []string) {
	for _, f := range arg {
		res = append(res, ``+filepath.FromSlash(f)+``)
	}
	return
}
