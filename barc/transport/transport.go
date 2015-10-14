package transport
import (
	"github.com/akaspin/bar/proto"
	"time"
	"github.com/akaspin/bar/manifest"
	"sync"
	"github.com/tamtam-im/logx"
	"fmt"
	"github.com/akaspin/bar/barc/lists"
	"github.com/akaspin/bar/barc/model"
	"bytes"
	"strings"
	"github.com/akaspin/bar/proto/bar"
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
}

func (t *Transport) Ping() (res proto.Info, err error) {
	cli, err := t.rpcPool.Take()
	if err != nil {
		return
	}
	defer cli.Release()

	err = cli.Call("Service.Ping", &struct{}{}, &res)
	return
}

// Upload blobs
func (t *Transport) Upload(blobs lists.Links) (err error) {
	// declare upload
	toUpload, err := t.NewUpload(blobs)
	if err != nil {
		return
	}

	// Ok. Now we can upload all
	// chunks before commit
	wg := sync.WaitGroup{}
	var errs []error

	for name, mt := range toUpload {
		wg.Add(1)
		go func(n string, mi manifest.Manifest) {
			defer wg.Done()
			var err1 error
			if err1 = t.UploadBlob(n, mi); err1 != nil {
				errs = append(errs, err1)
				return
			}
			if err1 = t.FinishUpload(mi.ID); err1 != nil {
				errs = append(errs, err1)
				return
			}

		}(name, mt)

	}
	wg.Wait()
	if len(errs) > 0 {
		err = fmt.Errorf("errors upload %s", errs)
		return
	}

	return
}

// Declare upload on bard server
func (t *Transport) NewUpload(blobs lists.Links) (toUpload lists.Links, err error) {
	var req, res []manifest.Manifest

	idmap := blobs.IDMap()

	for _, name := range idmap {
		req = append(req, blobs[name[0]])
	}

	cli, err := t.rpcPool.Take()
	if err != nil {
		return
	}
	defer cli.Release()

	if err = cli.Call("Service.NewUpload", &req, &res); err != nil {
		return
	}

	toUpload = lists.Links{}
	for _, m := range res {
		toUpload[idmap[m.ID][0]] = m
	}
	return
}

// Upload chunks from blob
func (t *Transport) UploadBlob(name string, info manifest.Manifest) (err error) {

	logx.Debugf("uploading %s %s (%d chunks)", name, info.ID, len(info.Chunks))

	var req, res []interface{}
	for _, v := range info.Chunks {
		req = append(req, v)
	}

	err = t.BatchPool.Do(
		func(ctx context.Context, in interface{}) (out interface{}, err error) {
			err = t.UploadChunk(name, info.ID, in.(manifest.Chunk))
			return
		}, &req, &res, concurrent.DefaultBatchOptions(),
	)
	return
}

//func (t *Transport) UploadBlob(name string, info manifest.Manifest) (err error) {
//
//	logx.Debugf("uploading %s %s (%d chunks)", name, info.ID, len(info.Chunks))
//
//	req := map[string]interface{}{}
//	for _, chunk := range info.Chunks {
//		req[chunk.ID.String()] = chunk
//	}
//
//	_, err = t.pool.RunBatch(parmap.Task{
//		Map: req,
//		Fn: func(id string, v interface{}) (res interface{}, err error) {
//			err = t.UploadChunk(name, info.ID, v.(manifest.Chunk))
//			return
//		},
//	})
//
//	return
//}

// Finish BLOB upload
func (t *Transport) FinishUpload(id manifest.ID) (err error) {
	cli, err := t.rpcPool.Take()
	if err != nil {
		return
	}
	defer cli.Release()

	var res struct{}
	if err = cli.Call("Service.CommitUpload", &id, &res); err != nil {
		return
	}
	logx.Debugf("finish upload BLOB %s", id)
	return
}

// Upload chunk
func (t *Transport) UploadChunk(name string, blobID manifest.ID, chunk manifest.Chunk) (err error) {
	logx.Tracef("uploading chunk %s (size: %d, offset %d) for BLOB %s:%s",
		chunk.ID, chunk.Size, chunk.Offset, name, blobID)

	buf := make([]byte, chunk.Size)
	if err = t.model.ReadChunk(name, chunk, buf); err != nil {
		return
	}

	cli, err := t.rpcPool.Take()
	if err != nil {
		return
	}
	defer cli.Release()

	var res struct{}
	err = cli.Call("Service.UploadChunk", &proto.ChunkData{
		proto.ChunkInfo{blobID, chunk}, buf,
	}, &res)

	logx.Tracef("chunk %s %s:%s %d uploaded", name,
		blobID, chunk.ID, chunk.Size)
	return
}

func (t *Transport) Download(blobs lists.Links) (err error) {

	logx.Tracef("fetching blobs %s", blobs.IDMap())
	fetch, err := t.GetManifests(blobs.IDMap().IDs())

	// little funny, but all chunks are equal - flatten them
	var req, res []interface{}
	chunkMap := map[string]struct{manifest.ID; manifest.Chunk}{}
	fetchIds := map[manifest.ID]struct{}{}

	for _, mt := range fetch {
		fetchIds[mt.ID] = struct{}{}
		for _, ch := range mt.Chunks    {
			chunkMap[ch.ID.String()] = struct{manifest.ID; manifest.Chunk}{mt.ID, ch}
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
			r := in.(struct{manifest.ID; manifest.Chunk})

			tclient, err := t.tPool.Take()
			if err != nil {
				return
			}
			defer tclient.Release()

			var blobID bar.ID
			if blobID, err = r.ID.MarshalThrift(); err != nil {
				return
			}
			var chunk bar.Chunk
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
	toAssemble := lists.Links{}
	for name, man := range blobs {
		_, exists := fetchIds[man.ID]
		if exists {
			toAssemble[name] = man
		}
	}

	err = a.Done(toAssemble)
	return
}

func (t *Transport) GetManifests(ids []manifest.ID) (res []manifest.Manifest, err error) {
	cli, err := t.tPool.Take()
	if err != nil {
		return
	}
	defer cli.Release()

	req, err := manifest.IDSlice(ids).MarshalThrift()
	if err != nil {
		return
	}

	res1, err := cli.GetManifests(req)

	var mx manifest.ManifestSlice
	err = (&mx).UnmarshalThrift(res1)
	res = mx
	return
}

func (t *Transport) Check(ids []manifest.ID) (res []manifest.ID, err error) {
	cli, err := t.rpcPool.Take()
	if err != nil {
		return
	}
	defer cli.Release()

	if err = cli.Call("Service.Check", &ids, &res); err != nil {
		return
	}

	return
}

func (t *Transport) UploadSpec(spec proto.Spec) (err error) {
	cli, err := t.rpcPool.Take()
	if err != nil {
		return
	}
	defer cli.Release()

	var res struct{}
	err = cli.Call("Service.UploadSpec", &spec, &res)

	return
}

func (t *Transport) GetSpec(id manifest.ID) (res lists.Links, err error) {
	cli, err := t.rpcPool.Take()
	if err != nil {
		return
	}
	defer cli.Release()
	err = cli.Call("Service.GetSpec", &id, &res)
	return
}