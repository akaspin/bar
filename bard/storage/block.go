package storage
import (
	"io"
	"path/filepath"
	"os"
	"fmt"
	"golang.org/x/crypto/sha3"
	"encoding/hex"
	"github.com/akaspin/go-contentaddressable"
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
	root string
	split int
}

func NewBlockStorage(root string, split int) *BlockStorage {
	return &BlockStorage{root, split}
}

func (s *BlockStorage) IsExists(id string) (ok bool, err error) {
	_, err = os.Stat(s.blobFileName(id))
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return
	}
	ok = true
	return

}

func (s *BlockStorage) StoreBLOB(id string, size int64, in io.Reader) (err error) {
	caOpts := contentaddressable.DefaultOptions()
	caOpts.Hasher = sha3.New256()
	w, err := contentaddressable.NewFileWithOptions(s.blobFileName(id), caOpts)
	if err != nil {
		return
	}
	defer w.Close()

	// TODO: squash internal hasher
	hasher := sha3.New256()

	mw := io.MultiWriter(w, hasher)

	written, err := io.Copy(mw, in)
	if err != nil {
		return
	}

	if written != size {
		err = fmt.Errorf("bad size for %s: actual %d != expected %d", id, written, size)
		return
	}
	if hex.EncodeToString(hasher.Sum(nil)) != id {
		err = fmt.Errorf("bad hash for %s not equal actual %s", id, hex.EncodeToString(hasher.Sum(nil)))
		return
	}
	err = w.Accept()
	return
}

func (s *BlockStorage) DestroyBLOB(id string) (err error) {
	err = os.Remove(s.blobFileName(id))
	return
}

func (s *BlockStorage) Close() (err error) {
	return
}

func (s *BlockStorage) blobFileName(id string) string {
	return filepath.Join(s.root, id[:s.split], id)
}
