package lists
import (
	"strings"
	"os"
	"path/filepath"
	"github.com/akaspin/bar/barc/git"
	"fmt"
)

type Globber struct {
	root string
	paths map[string]bool
	globs map[string]bool
	excludes map[string]bool
}

func NewGlobber(root string, in []string) (res *Globber) {
	res = &Globber{
		root,
		map[string]bool{},
		map[string]bool{},
		map[string]bool{},
	}
	for _, c := range in {
		if strings.HasPrefix(c, "!") {
			res.excludes[strings.TrimPrefix(c, "!")] = true
		} else if res.isGlob(c) {
			res.globs[c] = true
		} else {
			res.paths[c] = true
		}
	}
	if len(res.paths) == 0 {
		res.paths[""] = true
	}
	res.excludes[".git"] = true
	return
}

// List files in globber.
func (g *Globber) List(repo *git.Git) (res []string, err error) {
	for p, _ := range g.paths {
		res = append(res, g.walk(filepath.Join(g.root, p))...)
	}

	if repo != nil {
		res, err = repo.FilterByDiff("bar", res...)
		if err != nil {
			return
		}
		var dirty []string
		// fail if at least one file is dirty
		dirty, err = repo.DirtyFiles(res...)
		if err != nil {
			return
		}
		if len(dirty) > 0 {
			err = fmt.Errorf(
				"dirty files in working tree.\n\tgit add %s",
				strings.Join(dirty, " "))
		}
	}
	return
}

// Walk one path
func (g *Globber) walk(what string) (res []string) {
	info, err := os.Stat(what)
	if err != nil {
		return
	}

	if info.IsDir() {
		// walk
		err = filepath.Walk(what,
			func(path string, info os.FileInfo, inErr error) (outErr error) {
				if !info.IsDir() {
					res = append(res, g.walk(path)...)
				}
				return
			})
	} else {
		if g.checkIncludes(what) && !g.checkExcludes(what) {
			rel, err := filepath.Rel(g.root, what)
			if err == nil {
				res = append(res, rel)
			}
		}
	}

	return
}


func (g *Globber) checkIncludes(what string) (res bool) {
	if len(g.globs) == 0 {
		return true
	}
	rel, err := filepath.Rel(g.root, what)
	if err != nil {
		return true
	}
	for i, _ := range g.globs {
		ok, err := filepath.Match(i, rel)
		if err != nil {
			continue
		}
		if ok {
			return true
		}
	}
	return
}

func (g *Globber) checkExcludes(what string) (res bool) {
	if len(g.excludes) == 0 {
		return
	}
	rel, err := filepath.Rel(g.root, what)
	if err != nil {
		return
	}
	for i, _ := range g.excludes {
		if g.isGlob(i) {
			ok, err := filepath.Match(i, rel)
			if err != nil {
				continue
			}
			if ok {
				return true
			}
		} else {
			if strings.HasPrefix(rel, i) {
				return true
			}
		}

	}
	return
}

func (f *Globber) isGlob(s string) bool {
	return strings.ContainsAny(s, "*?[]^")
}

