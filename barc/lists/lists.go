package lists
import (
	"strings"
	"regexp"
	"path/filepath"
	"os"
)

type FileList struct {
	includes []string
	excludes []string
}

func NewFileList(arg ...string) (res *FileList) {
	res = &FileList{}
	for _, a := range arg {
		if strings.HasPrefix(a, "!") {
			res.excludes = append(res.excludes,
				normalizePattern(strings.TrimPrefix(a, "!")),
			)
		} else {
			res.includes = append(res.includes, normalizePattern(a))
		}
	}
	// add dotfiles to excludes
	res.excludes = append(res.excludes, "^\\..*$")
	res.excludes = append(res.excludes, "^.*/\\..*$")
	return
}

func (l *FileList) ListDir(dir string) (res []string) {
	filepath.Walk(dir, func(path string, info os.FileInfo, err1 error) error {
		if info.IsDir() {
			return nil
		}
		rel, err2 := filepath.Rel(dir, path)
		if err2 != nil {
			return nil
		}
		if l.isOK(filepath.ToSlash(rel)) {
			res = append(res, rel)
		}
		return nil
	})

	return
}


func (l *FileList) List(in []string) (res []string) {
	for _, name := range in {
		if l.isOK(name) {
			res = append(res, name)
		}
	}
	return
}

// what must be matched any of includes and none of excludes
func (l *FileList) isOK(what string) (res bool) {
	for _, ptn := range l.excludes {
		if ok, _ := regexp.MatchString(ptn, what); ok {
			return
		}
	}
	if len(l.includes) == 0 {
		res = true
		return
	}

	for _, ptn := range l.includes {
		if ok, _ := regexp.MatchString(ptn, what); ok {
			res = true
			return
		}
	}
	return
}

// replaces ** with .+ and * with [^/]+
func normalizePattern(in string) (out string) {
	out = strings.Replace(in, "**", ".+", -1)
	out = strings.Replace(out, "*", "[^/]+", -1)
	out = "^" + out + "$"
	return
}