package server
import (
	"github.com/akaspin/bar/proto/wire"
	"github.com/akaspin/bar/proto"
	"github.com/akaspin/bar/bard/storage"
	"bytes"
	"github.com/nu7hatch/gouuid"
	"time"
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

func (h *BardTHandler) CreateUpload(id []byte, manifests []*wire.Manifest, ttl int64) (r [][]byte, err error) {
	reqUploadID, err := uuid.Parse(id)
	if err != nil {
		return
	}
	var mans proto.ManifestSlice
	if err = (&mans).UnmarshalThrift(manifests); err != nil {
		return
	}
	r1, err := h.Storage.CreateUploadSession(*reqUploadID, mans, time.Duration(ttl))
	if err != nil {
		return
	}

	r, err = proto.IDSlice(r1).MarshalThrift()
	return
}

func (h *BardTHandler) UploadChunk(uploadId []byte, chunkId wire.ID, body []byte) (err error) {
	reqUploadID, err := uuid.Parse(uploadId)
	if err != nil {
		return
	}
	var id proto.ID
	if err = (&id).UnmarshalThrift(chunkId); err != nil {
		return
	}
	err = h.Storage.UploadChunk(*reqUploadID, id, bytes.NewReader(body))
	return
}

func (h *BardTHandler) FinishUpload(uploadId []byte) (err error) {
	reqUploadID, err := uuid.Parse(uploadId)
	if err != nil {
		return
	}
	err = h.Storage.FinishUploadSession(*reqUploadID)
	return
}

func (h *BardTHandler) GetMissingBlobIds(ids [][]byte) (r [][]byte, err error) {
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