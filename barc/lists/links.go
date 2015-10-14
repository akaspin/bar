package lists
import (
	"github.com/akaspin/bar/manifest"
	"path/filepath"
)

// Reverse mapping from id to names
type IDMap map[manifest.ID][]string

func (i IDMap) IDs() (res []manifest.ID) {
	for id, _ := range i {
		res = append(res, id)
	}
	return
}

type Links map[string]manifest.Manifest

func (l Links) IDMap() (res IDMap) {
	res = IDMap{}
	for name, m := range l {
		res[m.ID] = append(res[m.ID], name)
	}
	return
}

func (l Links) Names() (res []string) {
	for n, _ := range l {
		res = append(res, filepath.FromSlash(n))
	}
	return
}
