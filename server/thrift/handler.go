package thrift

import (
	"bytes"
	"fmt"
	"github.com/akaspin/bar/proto"
	"github.com/akaspin/bar/proto/wire"
	"github.com/akaspin/bar/server/storage"
	"github.com/nu7hatch/gouuid"
	"github.com/tamtam-im/logx"
	"time"
	"golang.org/x/net/context"
)

// Bar server thrift handler
type Handler struct {
	ctx context.Context
	*proto.ServerInfo
	storage.Storage
}

func NewHandler(ctx context.Context, info *proto.ServerInfo, stor storage.Storage) *Handler {
	return &Handler{ctx, info, stor}
}

func (h *Handler) GetInfo() (r *wire.ServerInfo, err error) {

	r1, err := h.ServerInfo.MarshalThrift()
	if err != nil {
		return
	}
	r = &r1
	return
}

func (h *Handler) CreateUpload(id []byte, manifests []*wire.Manifest, ttl int64) (r [][]byte, err error) {
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

func (h *Handler) UploadChunk(uploadId []byte, chunkId wire.ID, body []byte) (err error) {
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

func (h *Handler) FinishUpload(uploadId []byte) (err error) {
	reqUploadID, err := uuid.Parse(uploadId)
	if err != nil {
		return
	}
	if err = h.Storage.FinishUploadSession(*reqUploadID); err != nil {
		return
	}
	logx.Debugf("upload %s finished successfully", reqUploadID)
	return
}

func (h *Handler) GetMissingBlobIds(ids [][]byte) (r [][]byte, err error) {
	var req proto.IDSlice
	if err = (&req).UnmarshalThrift(ids); err != nil {
		return
	}

	res1, err := h.Storage.GetMissingBlobIDs(req)
	if err != nil {
		return
	}
	r, err = proto.IDSlice(res1).MarshalThrift()
	return
}

func (h *Handler) GetManifests(ids [][]byte) (r []*wire.Manifest, err error) {
	var req proto.IDSlice
	if err = (&req).UnmarshalThrift(ids); err != nil {
		return
	}
	logx.Debugf("serving manifests %s", req)

	res, err := h.Storage.GetManifests(req)
	if err != nil {
		logx.Error(err)
		return
	}

	r, err = proto.ManifestSlice(res).MarshalThrift()
	logx.Debugf("manifests served %s", req)

	return
}

func (h *Handler) FetchChunk(blobID wire.ID, chunk *wire.Chunk) (r []byte, err error) {
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

func (h *Handler) UploadSpec(spec *wire.Spec) (err error) {
	var req proto.Spec
	if err = (&req).UnmarshalThrift(*spec); err != nil {
		return
	}

	ok, err := h.Storage.IsSpecExists(req.ID)
	if err != nil {
		return
	}
	if ok {
		return
	}

	var ids []proto.ID
	for _, id := range req.BLOBs {
		ids = append(ids, id)
	}
	missing, err := h.Storage.GetMissingBlobIDs(ids)
	if err != nil {
		return
	}
	if len(missing) > 0 {
		err = fmt.Errorf("bad spec - missing BLOBs %s", missing)
	}

	err = h.Storage.WriteSpec(req)
	return
}

func (h *Handler) FetchSpec(id wire.ID) (r *wire.Spec, err error) {
	var req proto.ID
	if err = (&req).UnmarshalThrift(id); err != nil {
		return
	}
	res1, err := h.Storage.ReadSpec(req)
	if err != nil {
		return
	}
	res2, err := res1.MarshalThrift()
	if err != nil {
		return
	}
	r = &res2
	return
}
