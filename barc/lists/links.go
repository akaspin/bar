package lists
import (
	"github.com/akaspin/bar/proto"
	"path/filepath"
)

type Link struct {
	proto.Manifest
	Name string
}

// Reverse mapping from id to names
type IDMap map[proto.ID][]string

func (i IDMap) IDs() (res []proto.ID) {
	for id, _ := range i {
		res = append(res, id)
	}
	return
}

type Links map[string]proto.Manifest

func (l Links) ToSlice() (res []Link) {
	for k, v := range l {
		res = append(res, Link{v, k})
	}
	return
}

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
