package server
import (
	"github.com/akaspin/bar/proto/bar"
	"github.com/akaspin/bar/proto"
	"github.com/akaspin/bar/bard/storage"
)

type BardTHandler struct {
	*proto.Info
	storage.Storage
}

func NewBardTHandler(options *BardServerOptions) *BardTHandler {
	return &BardTHandler{options.Info, options.StoragePool}
}

func (h *BardTHandler) GetInfo() (r *bar.ServerInfo, err error) {
	r = &bar.ServerInfo{
		HttpEndpoint: h.Info.HTTPEndpoint,
		RpcEndpoints: h.Info.RPCEndpoints,
		ChunkSize: h.Info.ChunkSize,
		MaxConn: int32(h.Info.PoolSize),
		BufferSize: int32(h.Info.BufferSize),
	}
	return
}

func (h *BardTHandler) CreateUpload(id []byte, manifests []*bar.Manifest) (r []*bar.DataInfo, err error) {
	
	return
}

func (h *BardTHandler) UploadChunk(uploadId bar.ID, info *bar.DataInfo, body []byte) (err error) {
	return
}

func (h *BardTHandler) FinishUploadBlob(uploadId []byte, blobId bar.ID, tags [][]byte) (err error) {
	return
}

func (h *BardTHandler) FinishUpload(uploadId []byte) (err error) {
	return
}

func (h *BardTHandler) TagBlobs(ids [][]byte, tags [][]byte) (err error) {
	return
}

func (h *BardTHandler) UntagBlobs(ids [][]byte, tags [][]byte) (err error) {
	return
}

func (h *BardTHandler) IsBlobExists(ids [][]byte) (r [][]byte, err error) {
	return
}

func (h *BardTHandler) GetFetch(ids [][]byte) (r []*bar.Manifest, err error) {
	return
}

func (h *BardTHandler) FetchChunk(blobID bar.ID, chunkID bar.ID) (r []byte, err error) {
	return
}

func (h *BardTHandler) UploadSpec(spec *bar.Spec) (err error) {
	return
}

func (h *BardTHandler) FetchSpec(id bar.ID) (r *bar.Spec, err error) {
	return
}