package proto
import (
	"golang.org/x/crypto/sha3"
	"sort"
	"encoding/hex"
	"path/filepath"
)


// Tree spec
type Spec struct {

	// Spec ID is SHA3-256 hash of sorted filepath:manifest-id + kills-filepaths
	ID ID

	// Spec timestamp
	Timestamp int64

	// BLOB links
	BLOBs map[string]ID

	// Deleted filenames (not implemented)
	Remove []string
}

func NewSpec(timestamp int64, in map[string]ID, kill []string) (res Spec, err error) {
	hasher := sha3.New256()
	var idBuf []byte

	var names sort.StringSlice
	for n, _ := range in {
		names = append(names, n)
	}
	sort.Sort(names)

	kills := sort.StringSlice(kill)
	kills.Sort()

	drop := map[string]ID{}
	var removalsDrop []string

	for _, n := range names {
		if _, err = hasher.Write([]byte(filepath.ToSlash(n))); err != nil {
			return
		}
		if err = in[n].Decode(idBuf); err != nil {
			return
		}
		if _, err = hasher.Write(idBuf); err != nil {
			return
		}
		drop[filepath.ToSlash(n)] = in[n]
	}
	for _, n := range kills {
		if _, err = hasher.Write([]byte(filepath.ToSlash(n))); err != nil {
			return
		}
		removalsDrop = append(removalsDrop, filepath.ToSlash(n))
	}

	res = Spec{
		ID(hex.EncodeToString(hasher.Sum(nil))),
		timestamp,
		drop,
		removalsDrop}
	return
}