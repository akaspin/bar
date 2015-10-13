package transport
import (
	"github.com/akaspin/bar/proto"
	"time"
	"github.com/akaspin/bar/proto/manifest"
	"sync"
	"github.com/tamtam-im/logx"
	"fmt"
	"github.com/akaspin/bar/barc/lists"
	"github.com/akaspin/bar/barc/model"
	"github.com/akaspin/bar/parmap"
	"bytes"
	"strings"
	"encoding/hex"
	"github.com/akaspin/bar/proto/bar"
)

// Common transport with pooled connections
type Transport struct {
	model *model.Model
	rpcPool *RPCPool
	tPool *TPool
	pool *parmap.ParMap
}

// New RPC pool with default endpoint
func NewTransport(mod *model.Model, endpoint string, rpcEndpoints string, n int) (res *Transport) {
	rpcEP := strings.Split(rpcEndpoints, ",")
	res = &Transport{
		model: mod,
		rpcPool: NewRPCPool(n, time.Hour, endpoint, rpcEP),
		tPool: NewTPool(strings.Split(rpcEndpoints, ","), 1024 * 1024 * 8,  n, time.Hour),
		pool: parmap.NewWorkerPool(n),
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

	req := map[string]interface{}{}
	for _, chunk := range info.Chunks {
		req[chunk.ID] = chunk
	}

	_, err = t.pool.RunBatch(parmap.Task{
		Map: req,
		Fn: func(id string, v interface{}) (res interface{}, err error) {
			err = t.UploadChunk(name, info.ID, v.(manifest.Chunk))
			return
		},
	})

	return
}

// Finish BLOB upload
func (t *Transport) FinishUpload(id string) (err error) {
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
func (t *Transport) UploadChunk(name string, blobID string, chunk manifest.Chunk) (err error) {
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

	fetch, err := t.GetFetch(blobs.IDMap().IDs())
	logx.Debugf("fetching blobs %s", blobs.IDMap())

	// little funny, but all chunks are equal - flatten them
	chunkMap := map[string]interface{}{}
	fetchIds := map[string]struct{}{}

	for _, mt := range fetch {
		fetchIds[mt.ID] = struct{}{}
		for _, ch := range mt.Chunks    {
			chunkMap[ch.ID] = proto.ChunkInfo{mt.ID, ch}
		}
	}

	a, err := model.NewAssembler(t.model)
	defer a.Close()

	// Fetch all chunks
	_, err = t.model.Pool.RunBatch(parmap.Task{
		Map: chunkMap,
		Fn: func(id string, arg interface{}) (res interface{}, err error) {
//			cli, err := t.rpcPool.Take()
//			if err != nil {
//				return
//			}
//			defer cli.Release()

			tclient, err := t.tPool.Take()
			if err != nil {
				return
			}
			defer tclient.Release()

			ci := arg.(proto.ChunkInfo)

//			var data proto.ChunkData
//			if err = cli.Call("Service.FetchChunk", &ci, &data); err != nil {
//				return
//			}
			var blobId, chunkId []byte
			if blobId, err = hex.DecodeString(ci.BlobID); err != nil {
				return
			}
			if chunkId, err = hex.DecodeString(ci.ID); err != nil {
				return
			}
			chunk := &bar.Chunk{
				&bar.DataInfo{
					chunkId, ci.Size,
				},
				ci.Offset,
			}

			data, err := tclient.FetchChunk(blobId, chunk)
			if err != nil {
				return
			}

			err = a.StoreChunk(bytes.NewReader(data), ci.ID)
			return
		},
		IgnoreErrors: true,
	})
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

func (t *Transport) GetFetch(ids []string) (res []manifest.Manifest, err error) {
	cli, err := t.rpcPool.Take()
	if err != nil {
		return
	}
	defer cli.Release()
	err = cli.Call("Service.GetFetch", &ids, &res)
	return
}

func (t *Transport) Check(ids []string) (res []string, err error) {
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

func (t *Transport) GetSpec(id string) (res lists.Links, err error) {
	cli, err := t.rpcPool.Take()
	if err != nil {
		return
	}
	defer cli.Release()
	err = cli.Call("Service.GetSpec", &id, &res)
	return
}