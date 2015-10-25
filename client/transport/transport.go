package transport
import (
	"time"
	"github.com/akaspin/bar/proto"
	"github.com/tamtam-im/logx"
	"github.com/akaspin/bar/client/lists"
	"github.com/akaspin/bar/client/model"
	"bytes"
	"strings"
	"github.com/akaspin/bar/proto/wire"
	"golang.org/x/net/context"
	"github.com/akaspin/bar/concurrent"
)

// Common transport with pooled connections
type Transport struct {
	model *model.Model
	rpcPool *RPCPool
	tPool *TPool
	*concurrent.BatchPool
}

// New RPC pool with default endpoint
func NewTransport(mod *model.Model, endpoint string, rpcEndpoints string, n int) (res *Transport) {
	rpcEP := strings.Split(rpcEndpoints, ",")
	res = &Transport{
		model: mod,
		rpcPool: NewRPCPool(n, time.Hour, endpoint, rpcEP),
		tPool: NewTPool(strings.Split(rpcEndpoints, ","), 1024 * 1024 * 8,  n, time.Hour),
		BatchPool: concurrent.NewPool(n),
	}
	return
}

func (t *Transport) Close() {
	t.rpcPool.Close()
	t.tPool.Close()
	t.BatchPool.Close()
}

func (t *Transport) ServerInfo() (res proto.ServerInfo, err error) {
	cli, err := t.tPool.Take()
	if err != nil {
		return
	}
	defer cli.Release()

	res1, err := cli.GetInfo()

	err = (&res).UnmarshalThrift(*res1)
	return
}

// Upload blobs
func (t *Transport) Upload(blobs lists.BlobMap) (err error) {
	// declare upload
	upload := NewUpload(t, time.Hour)

	missing, err := upload.SendCreateUpload(blobs)
	if err != nil {
		return
	}

	var req, res []interface{}
	toUp := blobs.GetChunkLinkSlice(missing)
	for _, v := range toUp {
		req = append(req, v)
	}

	err = t.BatchPool.Do(
		func(ctx context.Context, in interface{}) (out interface{}, err error) {
			l := in.(lists.ChunkLink)
			err = upload.UploadChunk(l.Name, l.Chunk)
			return
		}, &req, &res, concurrent.DefaultBatchOptions().AllowErrors(),
	)
	logx.OnError(err)

	err = upload.Commit()
	return
}


func (t *Transport) Download(blobs lists.BlobMap) (err error) {

	logx.Tracef("fetching blobs %s", blobs.IDMap())
	fetch, err := t.GetManifests(blobs.IDMap().IDs())

	// little funny, but all chunks are equal - flatten them
	var req, res []interface{}
	chunkMap := map[string]struct{proto.ID; proto.Chunk}{}
	fetchIds := map[proto.ID]struct{}{}

	for _, mt := range fetch {
		fetchIds[mt.ID] = struct{}{}
		for _, ch := range mt.Chunks    {
			chunkMap[ch.ID.String()] = struct{proto.ID; proto.Chunk}{mt.ID, ch}
		}
	}
	for _, v := range chunkMap {
		req = append(req, v)

	}

	a, err := model.NewAssembler(t.model)
	defer a.Close()

	// Fetch all chunks
	err = t.model.BatchPool.Do(
		func (ctx context.Context, in interface{}) (out interface{}, err error) {
			r := in.(struct{proto.ID; proto.Chunk})

			tclient, err := t.tPool.Take()
			if err != nil {
				return
			}
			defer tclient.Release()

			var blobID wire.ID
			if blobID, err = r.ID.MarshalThrift(); err != nil {
				return
			}
			var chunk wire.Chunk
			chunk, err = r.Chunk.MarshalThrift()
			if err != nil {
				return
			}
			data, err := tclient.FetchChunk(blobID, &chunk)
			if err != nil {
				return
			}

			err = a.StoreChunk(bytes.NewReader(data), r.Chunk.ID)
			return
		}, &req, &res, concurrent.DefaultBatchOptions().AllowErrors(),
	)
	logx.OnError(err)

	// filter blobs
	toAssemble := lists.BlobMap{}
	for name, man := range blobs {
		_, exists := fetchIds[man.ID]
		if exists {
			toAssemble[name] = man
		}
	}

	err = a.Done(toAssemble)
	return
}

func (t *Transport) GetManifests(ids []proto.ID) (res []proto.Manifest, err error) {
	cli, err := t.tPool.Take()
	if err != nil {
		return
	}
	defer cli.Release()

	req, err := proto.IDSlice(ids).MarshalThrift()
	if err != nil {
		return
	}

	res1, err := cli.GetManifests(req)

	var mx proto.ManifestSlice
	err = (&mx).UnmarshalThrift(res1)
	res = mx
	return
}

func (t *Transport) Check(ids []proto.ID) (res []proto.ID, err error) {
	cli, err := t.tPool.Take()
	if err != nil {
		return
	}
	defer cli.Release()

	req, err := proto.IDSlice(ids).MarshalThrift()
	if err != nil {
		return
	}

	res1, err := cli.GetMissingBlobIds(req)
	if err != nil {
		return
	}

	var res2 proto.IDSlice
	err = (&res2).UnmarshalThrift(res1)
	res = res2

	return
}

func (t *Transport) UploadSpec(spec proto.Spec) (err error) {
	cli, err := t.tPool.Take()
	if err != nil {
		return
	}
	defer cli.Release()

	req, err := spec.MarshalThrift()
	if err != nil {
		return
	}
	err = cli.UploadSpec(&req)
	return
}

func (t *Transport) GetSpec(id proto.ID) (res proto.Spec, err error) {
	cli, err := t.tPool.Take()
	if err != nil {
		return
	}
	defer cli.Release()

	req, err := id.MarshalThrift()
	if err != nil {
		return
	}
	r1, err := cli.FetchSpec(req)
	if err != nil {
		return
	}
	err = (&res).UnmarshalThrift(*r1)

	return
}