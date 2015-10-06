package storage
import (
	"io"
	"path/filepath"
	"os"
	"fmt"
	"golang.org/x/crypto/sha3"
	"github.com/akaspin/go-contentaddressable"
	"github.com/akaspin/bar/proto/manifest"
	"encoding/json"
	"github.com/akaspin/bar/proto"
)

const (
	blob_ns = "blobs"
	spec_ns = "specs"
	manifests_ns = "manifests"
)


type BlockStorageFactory struct  {
	root string
	split int
}

func NewBlockStorageFactory(root string, split int) *BlockStorageFactory {
	return &BlockStorageFactory{root, split}
}

func (f *BlockStorageFactory) GetStorage() (StorageDriver, error)  {
	return NewBlockStorage(f.root, f.split), nil
}

// Simple block device storage
type BlockStorage struct {

	// Storage root
	Root string

	// Split factor
	Split int
}

func NewBlockStorage(root string, split int) *BlockStorage {
	return &BlockStorage{root, split}
}

func (s *BlockStorage) IsExists(id string) (ok bool, err error) {
	_, err = os.Stat(s.filePath(blob_ns, id))
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return
	}
	ok = true
	return
}

func (s *BlockStorage) WriteSpec(in proto.Spec) (err error) {
	w, err := os.OpenFile(s.filePath(spec_ns, in.ID),
		os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return
	}
	defer w.Close()
	err = json.NewEncoder(w).Encode(&in)
	return
}

func (s *BlockStorage) ReadSpec(id string) (res proto.Spec, err error) {
	r, err := os.Open(s.filePath(spec_ns, id))
	if err != nil {
		return
	}
	defer r.Close()
	res = proto.Spec{}
	err = json.NewDecoder(r).Decode(&res)
	return
}

func (s *BlockStorage) WriteManifest(m manifest.Manifest) (err error) {
	w, err := os.OpenFile(s.filePath(manifests_ns, m.ID),
		os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return
	}
	defer w.Close()
	err = json.NewEncoder(w).Encode(&m)
	return
}

func (s *BlockStorage) ReadManifest(id string) (res manifest.Manifest, err error) {
	r, err := os.Open(s.filePath(manifests_ns, id))
	if err != nil {
		return
	}
	defer r.Close()
	res = manifest.Manifest{}
	err = json.NewDecoder(r).Decode(&res)
	return
}

func (s *BlockStorage) WriteBLOB(id string, size int64, in io.Reader) (err error) {
	caOpts := contentaddressable.DefaultOptions()
	caOpts.Hasher = sha3.New256()
	w, err := contentaddressable.NewFileWithOptions(s.filePath(blob_ns, id), caOpts)
	if err != nil {
		return
	}
	defer w.Close()

	written, err := io.Copy(w, in)
	if err != nil {
		return
	}

	if written != size {
		err = fmt.Errorf("bad size for %s: actual %d != expected %d", id, written, size)
		return
	}
	err = w.Accept()
	return
}

func (s *BlockStorage) DestroyBLOB(id string) (err error) {
	err = os.Remove(s.filePath(blob_ns, id))
	return
}

func (s *BlockStorage) ReadBLOB(id string) (r io.ReadCloser, err error) {
	r, err = os.Open(s.filePath(blob_ns, id))
	return
}

func (s *BlockStorage) Close() (err error) {
	return
}

// Make filename
func (s *BlockStorage) filePath(what, id string) string {
	return filepath.Join(s.Root, what, id[:s.Split], id)
}
