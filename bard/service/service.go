package service
import (
	"github.com/akaspin/bar/proto"
	"github.com/akaspin/bar/bard/storage"
	"bytes"
	"github.com/tamtam-im/logx"
	"fmt"
	"github.com/akaspin/bar/barc/lists"
	"sync"
)

// RPC service
type Service struct {
	Info *proto.Info
	Storage storage.Storage
}

// Just returns server info
func (s *Service) Ping(req *struct{}, res *proto.Info) (err error) {
	*res = *s.Info
	logx.Debug("pong")
	return
}

// Check BLOBs
func (s *Service) Check(req *[]proto.ID, res *[]proto.ID) (err error) {


	logx.Debugf("checking %s", *req)

	var res1 []proto.ID
	for _, id := range *req {
		exists, err := s.Storage.IsBLOBExists(id)
		if err == nil && exists {
			res1 = append(res1, id)
		}
	}
	*res = res1
	return
}

// Takes manifests from client and returns missing BLOB ids
func (s *Service) NewUpload(req *[]proto.Manifest, res *[]proto.Manifest) (err error) {
	var missing []proto.Manifest
	var exists bool
	for _, m := range *req {
		if exists, err = s.Storage.IsBLOBExists(m.ID); err != nil {
			return
		}
		if !exists {
			if err = s.Storage.DeclareUpload(m); err != nil {
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
func (s *Service) CommitUpload(id *proto.ID, res *struct{}) (err error) {
	if err = s.Storage.FinishUpload(*id); err != nil {
		return
	}
	logx.Debugf("upload %s finished", *id)
	return
}

// Store chunk for declared blob
func (s *Service) UploadChunk(chunk *proto.ChunkData, res *struct{}) (err error) {

	err = s.Storage.WriteChunk(chunk.BlobID, chunk.Chunk.ID, chunk.Size,
		bytes.NewReader(chunk.Data))
	if err != nil {
		return
	}
	logx.Tracef("chunk stored %s:%s %d bytes",
		chunk.BlobID, chunk.Chunk.ID, chunk.Size)
	return
}

// Get manifests for download blobs
func (s *Service) GetFetch(req *[]proto.ID, res *[]proto.Manifest) (err error) {
	var feed []*proto.Manifest

	for _, id := range *req {
		m1, err := s.Storage.ReadManifest(id)
		if err != nil {
			return err
		}
		feed = append(feed, m1)
	}
	for _, v := range feed {
		*res = append(*res, *v)
	}

//	*res = *(&feed)
	return
}

// Upload spec
func (s *Service) UploadSpec(spec *proto.Spec, res *struct{}) (err error) {

	ok, err := s.Storage.IsSpecExists(spec.ID)
	if err != nil {
		return
	}
	if ok {
		return
	}

	for n, m := range spec.BLOBs {
		exists, err := s.Storage.IsBLOBExists(m)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("blob %s:%s not on bard for spec %s", n, m, spec.ID)
		}
	}

	if err = s.Storage.WriteSpec(*spec); err != nil {
		logx.Error(err)
		return
	}
	logx.Debugf("spec %s stored", spec.ID)
	return
}

// Get all links by spec-id
func (s *Service) GetSpec(id *proto.ID, res *lists.Links) (err error) {

	spec, err := s.Storage.ReadSpec(*id)
	if err != nil {
		return
	}

	res1 := lists.Links{}
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []error

	for name, manifestID := range spec.BLOBs {
		wg.Add(1)
		go func(n string, mID proto.ID) {
			defer wg.Done()

			var err1 error
			mu.Lock()
			rr, err1 := s.Storage.ReadManifest(mID)
			if err1 != nil {
				errs = append(errs, err1)
				return
			}
			res1[n] = *rr
			mu.Unlock()
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

