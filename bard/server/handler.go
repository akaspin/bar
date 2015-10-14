package server
import (
	"github.com/akaspin/bar/proto/bar"
	"github.com/akaspin/bar/proto"
	"github.com/akaspin/bar/bard/storage"
	"bytes"
	"github.com/akaspin/bar/manifest"
)

type BardTHandler struct {
	*proto.Info
	storage.Storage
}

func NewBardTHandler(options *BardServerOptions) *BardTHandler {
	return &BardTHandler{options.Info, options.Storage}
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

func (h *BardTHandler) GetManifests(ids [][]byte) (r []*bar.Manifest, err error) {
	var req []manifest.ID
	var res []manifest.Manifest

	for _, id := range ids {
		var i manifest.ID
		if err = (&i).UnmarshalThrift(id); err != nil {
			return
		}
		req = append(req, i)
	}

	res, err = h.Storage.GetManifests(req)
	if err != nil {
		return
	}

	r, err = manifest.ManifestSlice(res).MarshalThrift()
	return
}

func (h *BardTHandler) FetchChunk(blobID bar.ID, chunk *bar.Chunk) (r []byte, err error) {
	w := new(bytes.Buffer)
	var id manifest.ID
	if err = (&id).UnmarshalThrift(blobID); err != nil {
		return
	}
	if err = h.Storage.ReadChunkFromBlob(id, chunk.Info.Size, chunk.Offset, w); err != nil {
		return
	}
	r = w.Bytes()
	return
}

func (h *BardTHandler) UploadSpec(spec *bar.Spec) (err error) {
	return
}

func (h *BardTHandler) FetchSpec(id bar.ID) (r *bar.Spec, err error) {
	return
}