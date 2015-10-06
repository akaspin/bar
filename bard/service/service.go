package service
import (
	"github.com/akaspin/bar/proto"
	"github.com/akaspin/bar/bard/storage"
	"github.com/akaspin/bar/proto/manifest"
	"bytes"
	"github.com/tamtam-im/logx"
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
		if exists, err = store.IsExists(m.ID); err != nil {
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
func (s *Service) UploadChunk(chunk *proto.BLOBChunk, res *struct{}) (err error) {
	store, err := s.Storage.Take()
	if err != nil {
		return
	}
	defer s.Storage.Release(store)

	err = store.WriteChunk(chunk.BlobID, chunk.ChunkID, chunk.Size,
		bytes.NewReader(chunk.Data))
	if err != nil {
		return
	}
	logx.Debugf("chunk stored %s:%s %d bytes",
		chunk.BlobID, chunk.ChunkID, chunk.Size)
	return
}



