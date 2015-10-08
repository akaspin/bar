package service
import (
	"github.com/akaspin/bar/proto"
	"github.com/akaspin/bar/bard/storage"
	"github.com/akaspin/bar/proto/manifest"
	"bytes"
	"github.com/tamtam-im/logx"
	"fmt"
	"github.com/akaspin/bar/barc/lists"
	"sync"
)

// RPC service
type Service struct {
	Info *proto.Info
	Storage *storage.StoragePool
}

// Just returns server info
func (s *Service) Ping(req *struct{}, res *proto.Info) (err error) {
	*res = *s.Info
	return
}

// Check BLOBs
func (s *Service) Check(req *[]string, res *[]string) (err error) {
	store, err := s.Storage.Take()
	if err != nil {
		return
	}
	defer s.Storage.Release(store)

	logx.Debugf("checking %s", *req)

	var res1 []string
	for _, id := range *req {
		exists, err := store.IsBLOBExists(id)
		if err == nil && exists {
			res1 = append(res1, id)
		}
	}
	*res = res1
	return
}

// Takes manifests from client and returns missing BLOB ids
func (s *Service) NewUpload(req *[]manifest.Manifest, res *[]manifest.Manifest) (err error) {
	store, err := s.Storage.Take()
	if err != nil {
		return
	}
	defer s.Storage.Release(store)

	var missing []manifest.Manifest
	var exists bool
	for _, m := range *req {
		if exists, err = store.IsBLOBExists(m.ID); err != nil {
			return
		}
		if !exists {
			if err = store.DeclareUpload(m); err != nil {
				return
			}
			missing = append(missing, m)
			logx.Debugf("new upload %s declared", m.ID)
		}
	}
	*res = missing
	return
}

// Finish upload with ID
func (s *Service) CommitUpload(id *string, res *struct{}) (err error) {
	store, err := s.Storage.Take()
	if err != nil {
		return
	}
	defer s.Storage.Release(store)

	if err = store.FinishUpload(*id); err != nil {
		return
	}
	logx.Debugf("upload %s finished", *id)
	return
}

// Store chunk for declared blob
func (s *Service) UploadChunk(chunk *proto.ChunkData, res *struct{}) (err error) {
	store, err := s.Storage.Take()
	if err != nil {
		return
	}
	defer s.Storage.Release(store)

	err = store.WriteChunk(chunk.BlobID, chunk.Chunk.ID, chunk.Size,
		bytes.NewReader(chunk.Data))
	if err != nil {
		return
	}
	logx.Tracef("chunk stored %s:%s %d bytes",
		chunk.BlobID, chunk.Chunk.ID, chunk.Size)
	return
}

// Get manifests for download blobs
func (s *Service) GetFetch(req *[]string, res *[]manifest.Manifest) (err error) {
	var feed []manifest.Manifest
	store, err := s.Storage.Take()
	if err != nil {
		return
	}
	defer s.Storage.Release(store)

	for _, id := range *req {
		m1, err := store.ReadManifest(id)
		if err != nil {
			return err
		}
		feed = append(feed, m1)
	}
	*res = *(&feed)
	return
}

func (s *Service) FetchChunk(req *proto.ChunkInfo, res *proto.ChunkData) (err error) {
	store, err := s.Storage.Take()
	if err != nil {
		return
	}
	defer s.Storage.Release(store)

	buf := new(bytes.Buffer)
	err = store.ReadChunk(*req, buf)

	readed := &proto.ChunkData{*req, buf.Bytes()}
	*res = *readed

	return
}

// Upload spec
func (s *Service) UploadSpec(spec *proto.Spec, res *struct{}) (err error) {
	store, err := s.Storage.Take()
	if err != nil {
		return
	}
	defer s.Storage.Release(store)

	ok, err := store.IsSpecExists(spec.ID)
	if err != nil {
		return
	}
	if ok {
		return
	}

	for n, m := range spec.BLOBs {
		exists, err := store.IsBLOBExists(m)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("blob %s:%s not on bard for spec %s", n, m, spec.ID)
		}
	}

	if err = store.WriteSpec(*spec); err != nil {
		logx.Error(err)
		return
	}
	logx.Debugf("spec %s stored", spec.ID)
	return
}

// Get all links by spec-id
func (s *Service) GetSpec(id *string, res *lists.Links) (err error) {
	store, err := s.Storage.Take()
	if err != nil {
		return
	}
	defer s.Storage.Release(store)

	spec, err := store.ReadSpec(*id)
	if err != nil {
		return
	}

	res1 := lists.Links{}
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []error

	for name, manifestID := range spec.BLOBs {
		wg.Add(1)
		go func(n, mID string) {
			defer wg.Done()
			store1, err := s.Storage.Take()
			if err != nil {
				return
			}
			defer s.Storage.Release(store1)

			var err1 error
			mu.Lock()
			res1[n], err1 = store1.ReadManifest(mID)
			mu.Unlock()
			if err1 != nil {
				errs = append(errs, err1)
				return
			}
		}(name, manifestID)
	}
	wg.Wait()

	if len(errs) > 0 {
		err = fmt.Errorf("errors while collecting manifests for spec %s: %s",
			spec.ID, errs)
		return
	}

	*res = res1
	return
}

