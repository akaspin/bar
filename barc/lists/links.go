package lists
import (
	"github.com/akaspin/bar/proto/manifest"
	"path/filepath"
)

type IDMap map[string]string

func (i IDMap) IDs() (res []string) {
	for id, _ := range i {
		res = append(res, id)
	}
	return
}

type Links map[string]manifest.Manifest

func (l Links) IDMap() (res IDMap) {
	res = IDMap{}
	for name, m := range l {
		res[m.ID] = name
	}
	return
}

func (l Links) Names() (res []string) {
	for n, _ := range l {
		res = append(res, filepath.FromSlash(n))
	}
	return
}

func (l Links) FromIDMap(idmap IDMap) (res Links) {
	res = Links{}
	for _, name := range idmap {
		res[name] = l[name]
	}
	return
}