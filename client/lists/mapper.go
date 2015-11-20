package lists

import (
	"path/filepath"
	"strings"
	"github.com/tamtam-im/logx"
)

func OSFromSlash(name string) (res string)  {
	res = filepath.Join(filepath.SplitList(filepath.FromSlash(name))...)
	return
}

func OSJoin(chunks ...string) (res string) {
	var osChunks []string
	var escaped bool
	for _, c := range chunks {
		if strings.HasPrefix(c, `"`) || strings.HasSuffix(c, `"`) {
			c = strings.Trim(c, `"`)
			escaped = true
		}
		osChunks = append(osChunks, filepath.Clean(filepath.FromSlash(c)))
	}
	res = filepath.ToSlash(filepath.Join(osChunks...))
	if escaped {
		res = `"` + res + `"`
	}
	return
}

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
	logx.Tracef("remapping from %s to %s %s", from, to, arg)
	if filepath.ToSlash(from) == filepath.ToSlash(to) {
		res = arg
		return
	}
	var one string
	for _, p := range arg {
		if one, err = filepath.Rel(OSFromSlash(to), OSFromSlash(OSJoin(from, p))); err != nil {
			logx.Errorf(">>> %s %s %s", err, OSJoin(from, p), to)
			return
		}
		res = append(res, one)
	}
	logx.Tracef("remapped %s", res)

	return
}

func (m *Mapper) ToShell(arg ...string) (res []string) {
	for _, f := range arg {
		res = append(res, filepath.FromSlash(f))
	}
	return
}
