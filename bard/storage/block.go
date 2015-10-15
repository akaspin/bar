package storage
import (
	"io"
	"path/filepath"
	"os"
	"fmt"
	"golang.org/x/crypto/sha3"
	"github.com/akaspin/go-contentaddressable"
	"github.com/akaspin/bar/proto"
	"encoding/json"
	"time"
	"github.com/akaspin/bar/concurrent"
	"golang.org/x/net/context"
)

const (
	blob_ns = "blobs"
	spec_ns = "specs"
	manifests_ns = "manifests"
	upload_ns = "uploads"
)

type BlockStorageOptions struct {
	// Storage root
	Root string

	// Split factor
	Split int

	MaxFiles int
	PoolSize int
}


// Simple block device storage
type BlockStorage struct {

	*BlockStorageOptions

	// Max Open files locker
	FDLocks *concurrent.LocksPool

	*concurrent.BatchPool
}

func NewBlockStorage(options *BlockStorageOptions) *BlockStorage {
	return &BlockStorage{
		BlockStorageOptions: options,
		FDLocks: concurrent.NewLockPool(options.MaxFiles, time.Minute * 5),
		BatchPool: concurrent.NewPool(options.PoolSize),
	}
}

func (s *BlockStorage) IsSpecExists(id proto.ID) (ok bool, err error) {
	ok, err = s.isExists(spec_ns, proto.ID(id.String() + ".json"))
	return
}

func (s *BlockStorage) IsBLOBExists(id proto.ID) (ok bool, err error) {
	ok, err = s.isExists(blob_ns, id)
	return
}

func (s *BlockStorage) isExists(ns string, id proto.ID) (ok bool, err error) {
	lock, err := s.FDLocks.Take()
	if err != nil {
		return
	}
	defer lock.Close()

	_, err = os.Stat(s.idPath(ns, id))
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
	lock, err := s.FDLocks.Take()
	if err != nil {
		return
	}
	defer lock.Close()

	specName := s.idPath(spec_ns, in.ID) + ".json"
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

func (s *BlockStorage) ReadSpec(id proto.ID) (res proto.Spec, err error) {
	lock, err := s.FDLocks.Take()
	if err != nil {
		return
	}
	defer lock.Close()

	r, err := os.Open(s.idPath(spec_ns, id) + ".json")
	if err != nil {
		return
	}
	defer r.Close()
	res = proto.Spec{}
	err = json.NewDecoder(r).Decode(&res)
	return
}

func (s *BlockStorage) ReadManifest(id proto.ID) (res *proto.Manifest, err error) {
	lock, err := s.FDLocks.Take()
	if err != nil {
		return
	}
	defer lock.Close()

	r, err := os.Open(s.idPath(manifests_ns, id))
	if err != nil {
		return
	}
	defer r.Close()
	res = &proto.Manifest{}
	err = json.NewDecoder(r).Decode(&res)
	return
}

func (s *BlockStorage) GetManifests(ids []proto.ID) (res []proto.Manifest, err error) {
	var req, res1 []interface{}
	for _, i := range ids {
		req = append(req, i)
	}

	if err = s.BatchPool.Do(
		func(ctx context.Context, in interface{}) (out interface{}, err error) {
			r := in.(proto.ID)
			r2, err := s.ReadManifest(r)
			if err != nil {
				return
			}
			out = *r2
			return
		}, &req, &res1, concurrent.DefaultBatchOptions(),
	); err != nil {
		return
	}

	for _, v := range res1 {
		res = append(res, v.(proto.Manifest))
	}

	return
}

func (s *BlockStorage) DeclareUpload(m proto.Manifest) (err error) {
	lock, err := s.FDLocks.Take()
	if err != nil {
		return
	}
	defer lock.Close()

	base := s.idPath(upload_ns, m.ID)
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

func (s *BlockStorage) WriteChunk(blobID, chunkID proto.ID, size int64, r io.Reader) (err error) {
	lock, err := s.FDLocks.Take()
	if err != nil {
		return
	}
	defer lock.Close()

	n := filepath.Join(s.idPath(upload_ns, blobID), chunkID.String())
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

func (s *BlockStorage) FinishUpload(id proto.ID) (err error) {

	base := s.idPath(upload_ns, id)
	manifestName := filepath.Join(base, "manifest.json")

	mr, err := os.Open(manifestName)
	if err != nil {
		return
	}
	defer mr.Close()

	var m proto.Manifest
	err = json.NewDecoder(mr).Decode(&m)
	if err != nil {
		return
	}

	w, err := s.getCAFile(s.idPath(blob_ns, id))
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
	targetManifest := s.idPath(manifests_ns, id)
	os.MkdirAll(filepath.Dir(targetManifest), 0755)
	err = os.Rename(manifestName, targetManifest)

	defer os.RemoveAll(base)
	return
}

func (s *BlockStorage) finishChunk(base string, info proto.Chunk, w io.Writer) (err error) {
	name := filepath.Join(base, string(info.ID))
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
	lock, err := s.FDLocks.Take()
	if err != nil {
		return
	}
	defer lock.Close()

	f, err := os.Open(s.idPath(blob_ns, chunk.BlobID))
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

func (s *BlockStorage) ReadChunkFromBlob(blobID proto.ID, size, offset int64, w io.Writer) (err error) {
	lock, err := s.FDLocks.Take()
	if err != nil {
		return
	}
	defer lock.Close()

	f, err := os.Open(s.idPath(blob_ns, blobID))
	if err != nil {
		return
	}
	defer f.Close()

	if _, err = f.Seek(offset, 0); err != nil {
		return
	}

	written, err := io.CopyN(w, f, size)
	if err != nil {
		return
	}
	if written != size {
		err = fmt.Errorf("bad size for chunk %s (offset %d) %d != %d",
			blobID, offset, size, written)
	}

	return
}

func (s *BlockStorage) Close() (err error) {
	return
}

func (s *BlockStorage) idPath(ns string, id proto.ID) string {
	ids := id.String()
	return filepath.Join(s.Root, ns, ids[:s.Split], ids)
}

func (s *BlockStorage) getCAFile(name string) (w *contentaddressable.File, err error) {
	caOpts := contentaddressable.DefaultOptions()
	caOpts.Hasher = sha3.New256()

	w, err = contentaddressable.NewFileWithOptions(name, caOpts)
	return
}
