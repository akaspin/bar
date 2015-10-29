package lists

import (
	"github.com/akaspin/bar/proto"
	"path/filepath"
)

// Link from filename to manifest
type BlobLink struct {
	proto.Manifest
	Name string
}

// Link to chunk in blob
type ChunkLink struct {
	Name string
	proto.Chunk
}

// Reverse mapping from id to names
type IDMap map[proto.ID][]string

func (i IDMap) ToBlobMap(manifests []proto.Manifest) (res BlobMap) {
	res = BlobMap{}
	for _, manifest := range manifests {
		if names, ok := i[manifest.ID]; ok {
			for _, name := range names {
				res[name] = manifest
			}
		}
	}
	return
}

func (i IDMap) IDs() (res []proto.ID) {
	for id, _ := range i {
		res = append(res, id)
	}
	return
}

// filename to manifest mapping
type BlobMap map[string]proto.Manifest

// Get unique chunk links
func (l BlobMap) GetChunkLinkSlice(chunkIDs []proto.ID) (res []ChunkLink) {
	// make chunk ref
	ref := map[proto.ID]struct{}{}
	for _, v := range chunkIDs {
		ref[v] = struct{}{}
	}

	var ok bool
	for name, man := range l {
		for _, chunk := range man.Chunks {
			_, ok = ref[chunk.ID]
			if ok {
				res = append(res, ChunkLink{name, chunk})
				delete(ref, chunk.ID)
				if len(ref) == 0 {
					return
				}
			}
		}
	}

	return
}

// Get unique manifests
func (l BlobMap) GetManifestSlice() (res []proto.Manifest) {
	ref := map[proto.ID]proto.Manifest{}
	for _, m := range l {
		ref[m.ID] = m
	}

	for _, v := range ref {
		res = append(res, v)
	}
	return
}

func (l BlobMap) ToSlice() (res []BlobLink) {
	for k, v := range l {
		res = append(res, BlobLink{v, k})
	}
	return
}

func (l BlobMap) IDMap() (res IDMap) {
	res = IDMap{}
	for name, m := range l {
		res[m.ID] = append(res[m.ID], name)
	}
	return
}

func (l BlobMap) Names() (res []string) {
	for n, _ := range l {
		res = append(res, filepath.FromSlash(n))
	}
	return
}
