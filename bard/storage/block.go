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
	upload_ns = "uploads"
)


type BlockStorageFactory struct {
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

func (s *BlockStorage) IsSpecExists(id string) (ok bool, err error) {
	ok, err = s.isExists(spec_ns, id + ".json")
	return
}

func (s *BlockStorage) IsBLOBExists(id string) (ok bool, err error) {
	ok, err = s.isExists(blob_ns, id)
	return
}

func (s *BlockStorage) isExists(ns string, id string) (ok bool, err error) {
	_, err = os.Stat(s.filePath(ns, id))
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
	specName := s.filePath(spec_ns, in.ID + ".json")
	if err = os.MkdirAll(filepath.Dir(specName), 0755); err != nil {
		return
	}
	w, err := os.OpenFile(specName,
		os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return
	}
	defer w.Close()
	err = json.NewEncoder(w).Encode(&in)
	return
}

func (s *BlockStorage) ReadSpec(id string) (res proto.Spec, err error) {
	r, err := os.Open(s.filePath(spec_ns, id + ".json"))
	if err != nil {
		return
	}
	defer r.Close()
	res = proto.Spec{}
	err = json.NewDecoder(r).Decode(&res)
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

func (s *BlockStorage) DeclareUpload(m manifest.Manifest) (err error) {
	base := s.filePath(upload_ns, m.ID)
	if err = os.MkdirAll(base, 0755); err != nil {
		return
	}

	w, err := os.OpenFile(filepath.Join(base, "manifest.json"),
		os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return
	}
	defer w.Close()
	err = json.NewEncoder(w).Encode(&m)
	return
}

func (s *BlockStorage) WriteChunk(blobID, chunkID string, size int64, r io.Reader) (err error) {
	n := filepath.Join(s.filePath(upload_ns, blobID), chunkID)
	w, err := s.getCAFile(n)
	if err != nil {
		return
	}
	defer w.Close()

	written, err := io.Copy(w, r)
	if err != nil {
		return
	}

	if written != size {
		err = fmt.Errorf("bad chunk size for %s:%s : %d != %d",
			blobID, chunkID, size, written)
		return
	}
	err = w.Accept()
	return
}

func (s *BlockStorage) FinishUpload(id string) (err error) {

	base := s.filePath(upload_ns, id)
	manifestName := filepath.Join(base, "manifest.json")

	mr, err := os.Open(manifestName)
	if err != nil {
		return
	}
	defer mr.Close()

	var m manifest.Manifest
	err = json.NewDecoder(mr).Decode(&m)
	if err != nil {
		return
	}

	w, err := s.getCAFile(s.filePath(blob_ns, id))
	if err != nil {
		return
	}
	defer w.Close()

	for _, chunkInfo := range m.Chunks {
		if err = s.finishChunk(base, chunkInfo, w); err != nil {
			return
		}
	}
	if err = w.Accept(); err != nil {
		return
	}

	// Relocate manifest
	targetManifest := s.filePath(manifests_ns, id)
	os.MkdirAll(filepath.Dir(targetManifest), 0755)
	err = os.Rename(manifestName, targetManifest)

	defer os.RemoveAll(base)
	return
}

func (s *BlockStorage) finishChunk(base string, info manifest.Chunk, w io.Writer) (err error) {
	name := filepath.Join(base, info.ID)
	r, err := os.Open(name)
	if err != nil {
		return
	}
	defer r.Close()

	written, err := io.Copy(w, r)
	if err != nil {
		return
	}
	if written != info.Size {
		err = fmt.Errorf("bad size for chunk %s %d != %s", name, info.Size, written)
	}

	return
}

func (s *BlockStorage) ReadChunk(chunk proto.ChunkInfo, w io.Writer) (err error) {
	f, err := os.Open(s.filePath(blob_ns, chunk.BlobID))
	if err != nil {
		return
	}
	defer f.Close()

	if _, err = f.Seek(chunk.Offset, 0); err != nil {
		return
	}

	written, err := io.CopyN(w, f, chunk.Size)
	if err != nil {
		return
	}
	if written != chunk.Size {
		err = fmt.Errorf("bad size for chunk %s %d != %s", chunk.ID, chunk.Size, written)
	}
	return
}

func (s *BlockStorage) DestroyBLOB(id string) (err error) {
	err = os.Remove(s.filePath(blob_ns, id))
	return
}

func (s *BlockStorage) Close() (err error) {
	return
}

// Make filename
func (s *BlockStorage) filePath(what, id string) string {
	return filepath.Join(s.Root, what, id[:s.Split], id)
}

func (s *BlockStorage) getCAFile(name string) (w *contentaddressable.File, err error) {
	caOpts := contentaddressable.DefaultOptions()
	caOpts.Hasher = sha3.New256()

	w, err = contentaddressable.NewFileWithOptions(name, caOpts)
	return
}
