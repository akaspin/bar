package server
import (
	"github.com/akaspin/bar/proto/wire"
	"github.com/akaspin/bar/proto"
	"github.com/akaspin/bar/bard/storage"
	"bytes"
)

type BardTHandler struct {
	*proto.Info
	storage.Storage
}

func NewBardTHandler(options *BardServerOptions) *BardTHandler {
	return &BardTHandler{options.Info, options.Storage}
}

func (h *BardTHandler) GetInfo() (r *wire.ServerInfo, err error) {
	r = &wire.ServerInfo{
		HttpEndpoint: h.Info.HTTPEndpoint,
		RpcEndpoints: h.Info.RPCEndpoints,
		ChunkSize: h.Info.ChunkSize,
		MaxConn: int32(h.Info.PoolSize),
		BufferSize: int32(h.Info.BufferSize),
	}
	return
}

func (h *BardTHandler) CreateUpload(id []byte, manifests []*wire.Manifest) (r []*wire.DataInfo, err error) {
	
	return
}

func (h *BardTHandler) UploadChunk(uploadId wire.ID, info *wire.DataInfo, body []byte) (err error) {
	return
}

func (h *BardTHandler) FinishUploadBlob(uploadId []byte, blobId wire.ID, tags [][]byte) (err error) {
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

func (h *BardTHandler) GetManifests(ids [][]byte) (r []*wire.Manifest, err error) {
	var req proto.IDSlice
	if err = (&req).UnmarshalThrift(ids); err != nil {
		return
	}

	res, err := h.Storage.GetManifests(req)
	if err != nil {
		return
	}

	r, err = proto.ManifestSlice(res).MarshalThrift()
	return
}

func (h *BardTHandler) FetchChunk(blobID wire.ID, chunk *wire.Chunk) (r []byte, err error) {
	w := new(bytes.Buffer)
	var id proto.ID
	if err = (&id).UnmarshalThrift(blobID); err != nil {
		return
	}
	if err = h.Storage.ReadChunkFromBlob(id, chunk.Info.Size, chunk.Offset, w); err != nil {
		return
	}
	r = w.Bytes()
	return
}

func (h *BardTHandler) UploadSpec(spec *wire.Spec) (err error) {

	return
}

func (h *BardTHandler) FetchSpec(id wire.ID) (r *wire.Spec, err error) {
	return
}