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
	"github.com/akaspin/concurrency"
	"golang.org/x/net/context"
	"github.com/nu7hatch/gouuid"
	"strings"
	"encoding/hex"
	"github.com/tamtam-im/logx"
)

const (
	blob_ns = "blobs"
	spec_ns = "specs"
	manifests_ns = "manifests"
	upload_ns = "uploads"
)

func BlockStorageFactory(opts map[string]string) (res Storage, err error) {
	op := &BlockStorageOptions{"data", 2, 32, 32}
	var ok bool
	var tmp string

	if tmp, ok = opts["root"]; ok {
		op.Root = tmp
	}
	if tmp, ok = opts["split"]; ok {
		if _, err = fmt.Sscanf(tmp, "%d", op.Split); err != nil {
			return
		}
	}
	if tmp, ok = opts["max-files"]; ok {
		if _, err = fmt.Sscanf(tmp, "%d", op.MaxFiles); err != nil {
			return
		}
	}
	if tmp, ok = opts["pool"]; ok {
		if _, err = fmt.Sscanf(tmp, "%d", op.PoolSize); err != nil {
			return
		}
	}

	res = NewBlockStorage(op)
	return
}

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
	FDLocks *concurrency.Locks

	*concurrency.BatchPool
}

func NewBlockStorage(options *BlockStorageOptions) *BlockStorage {
	return &BlockStorage{
		BlockStorageOptions: options,
		FDLocks: concurrency.NewLocks(context.Background(), options.MaxFiles, time.Minute * 5),
		BatchPool: concurrency.NewPool(options.PoolSize),
	}
}

func (s *BlockStorage) IsSpecExists(id proto.ID) (ok bool, err error) {
	err = s.FDLocks.With(func() (err error) {
		_, err = os.Stat(s.idPath(spec_ns, id) + ".json")
		if os.IsNotExist(err) {
			return nil
		}
		if err != nil {
			return
		}
		ok = true
		return
	})
	return
}

func (s *BlockStorage) WriteSpec(in proto.Spec) (err error) {
	lock, err := s.FDLocks.Take()
	if err != nil {
		return
	}
	defer lock.Release()

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
	defer lock.Release()

	r, err := os.Open(s.idPath(spec_ns, id) + ".json")
	if err != nil {
		return
	}
	defer r.Close()
	res = proto.Spec{}
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
			out, err = s.readManifest(s.idPath(manifests_ns, r) + ".json")
			return
		}, &req, &res1, concurrency.DefaultBatchOptions(),
	); err != nil {
		return
	}

	for _, v := range res1 {
		res = append(res, v.(proto.Manifest))
	}

	return
}

func (s *BlockStorage) GetMissingBlobIDs(ids []proto.ID) (res []proto.ID, err error) {
	lock, err := s.FDLocks.Take()
	if err != nil {
		return
	}
	defer lock.Release()

	for _, id := range ids {
		_, err = os.Stat(s.idPath(manifests_ns, id) + ".json")
		if os.IsNotExist(err) {
			res = append(res, id)
			err = nil
		} else if err != nil {
			return
		}
	}

	return
}


func (s *BlockStorage) WriteChunk(blobID, chunkID proto.ID, size int64, r io.Reader) (err error) {
	lock, err := s.FDLocks.Take()
	if err != nil {
		return
	}
	defer lock.Release()

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

func (s *BlockStorage) ReadChunkFromBlob(blobID proto.ID, size, offset int64, w io.Writer) (err error) {
	lock, err := s.FDLocks.Take()
	if err != nil {
		return
	}
	defer lock.Release()

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


func (s *BlockStorage) CreateUploadSession(uploadID uuid.UUID, in []proto.Manifest, ttl time.Duration) (missing []proto.ID, err error) {
	hexid := proto.ID(hex.EncodeToString(uploadID[:]))

	// take lock
	lock, err := s.FDLocks.Take()
	if err != nil {
		return
	}
	defer lock.Release()

	// Create directories and write support data
	base := filepath.Join(s.idPath(upload_ns,
		proto.ID(hex.EncodeToString(uploadID[:]))), manifests_ns)
	if err = os.MkdirAll(base, 0755); err != nil {
		return
	}

	var missingBlobs []proto.Manifest

	for _, m := range in {
		if err = func(m proto.Manifest) (err error) {
			var statErr error
			_, statErr = os.Stat(s.idPath(manifests_ns, m.ID))
			if os.IsNotExist(statErr) {
				missingBlobs = append(missingBlobs, m)
			} else if statErr != nil {
				err = statErr
				return
			} else {
				// exists - ok
				return
			}

			w, err := os.OpenFile(filepath.Join(base, m.ID.String() + "-manifest.json"),
				os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
			if err != nil {
				return
			}
			defer w.Close()
			err = json.NewEncoder(w).Encode(&m)
			return
		}(m); err != nil {
			return
		}
	}

	missing = proto.ManifestSlice(missingBlobs).GetChunkSlice()

	w, err := os.OpenFile(filepath.Join(s.idPath(upload_ns, hexid), "expires.timestamp"),
		os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return
	}
	defer w.Close()

	if _, err = w.Write([]byte(fmt.Sprintf("%d", time.Now().Add(ttl).UnixNano()))); err != nil {
		return
	}
	logx.Debugf("upload session %s created succefully", hexid)
	return
}

func (s *BlockStorage) UploadChunk(uploadID uuid.UUID, chunkID proto.ID, r io.Reader) (err error) {
	lock, err := s.FDLocks.Take()
	if err != nil {
		return
	}
	defer lock.Release()
	hexid := proto.ID(hex.EncodeToString(uploadID[:]))

	n := filepath.Join(s.idPath(upload_ns, hexid), chunkID.String())
	w, err := s.getCAFile(n)
	if err != nil {
		return
	}
	defer w.Close()

	if _, err = io.Copy(w, r); err != nil {
		return
	}

	err = w.Accept()
	return
}

func (s *BlockStorage) FinishUploadSession(uploadID uuid.UUID) (err error) {
	hexid := proto.ID(hex.EncodeToString(uploadID[:]))
	base := s.idPath(upload_ns, hexid)
	defer os.RemoveAll(base)


	// load manifests
	manifests_base := filepath.Join(base, manifests_ns)
	var manifests []proto.Manifest
	if err = func () (err error) {
		lock, err := s.FDLocks.Take()
		if err != nil {
			return
		}
		defer lock.Release()

		err = filepath.Walk(manifests_base, func (path string, info os.FileInfo, ferr error) (err error){
			if strings.HasSuffix(path, "-manifest.json") {
				var man proto.Manifest
				if man, err = s.readManifest(path); err != nil {
					return
				}
				manifests = append(manifests, man)
			}
			return
		})
		return
	}(); err != nil {
		return
	}

	// collect all manifests
	var req, res []interface{}
	for _, v := range manifests {
		req = append(req, v)
	}

	err = s.BatchPool.Do(
		func(ctx context.Context, in interface{}) (out interface{}, err error) {
			lock, err := s.FDLocks.Take()
			if err != nil {
				return
			}
			defer lock.Release()

			m := in.(proto.Manifest)
			target := s.idPath(blob_ns, m.ID)

			f, fErr := s.getCAFile(target)
			if os.IsExist(fErr) {
				return
			} else if fErr != nil {
				err = fErr
				return
			}
			defer f.Close()
			logx.Debugf("assembling %s", m.ID)

			for _, chunk := range m.Chunks {
				if err = func(chunk proto.Chunk) (err error) {
					lock, err := s.FDLocks.Take()
					if err != nil {
						return
					}
					defer lock.Release()

					r, err := os.Open(filepath.Join(base, chunk.ID.String()))
					if err != nil {
						return
					}
					defer r.Close()

					_, err = io.Copy(f, r)
					return
				}(chunk); err != nil {
					return
				}
			}
			err = f.Accept()

			// move manifest
			manTarget := s.idPath(manifests_ns, m.ID) + ".json"
			os.MkdirAll(filepath.Dir(manTarget), 0755)
			err = os.Rename(filepath.Join(manifests_base, m.ID.String() + "-manifest.json"), manTarget)
			return
		}, &req, &res, concurrency.DefaultBatchOptions().AllowErrors(),
	)
	return
}

func (s *BlockStorage) readManifest(filename string) (res proto.Manifest, err error) {
	lock, err := s.FDLocks.Take()
	if err != nil {
		return
	}
	defer lock.Release()

	r, err := os.Open(filename)
	if err != nil {
		return
	}
	defer r.Close()
	err = json.NewDecoder(r).Decode(&res)
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
