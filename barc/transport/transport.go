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
)

// Common transport with pooled connections
type Transport struct {
	WD string

	// base endpoint. http://example.com/v1
	DefaultEndpoint string

	model *model.Model
	rpcPool *RPCPool
	pool *parmap.ParMap
}

// New RPC pool with default endpoint
func NewTransport(mod *model.Model, endpoint string, n int) (res *Transport) {
	res = &Transport{
		WD: mod.WD,
		model: mod,
		DefaultEndpoint: endpoint,
		rpcPool: NewRPCPool(n, time.Hour, endpoint),
		pool: parmap.NewWorkerPool(n),
	}
	return
}

func (t *Transport) Close() {
	t.rpcPool.Close()
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
		req = append(req, blobs[name])
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
		toUpload[idmap[m.ID]] = m
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
	logx.Debug("fetching blobs %s", blobs.IDMap())

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
			cli, err := t.rpcPool.Take()
			if err != nil {
				return
			}
			defer cli.Release()

			ci := arg.(proto.ChunkInfo)

			var data proto.ChunkData
			if err = cli.Call("Service.FetchChunk", &ci, &data); err != nil {
				return
			}

			err = a.StoreChunk(bytes.NewReader(data.Data), ci.ID)
			return
		},
	})
	if err != nil {
		return
	}

	// filter blobs
	toAssemble := lists.Links{}
	for name, man := range blobs {
		_, exists := fetchIds[man.ID]
		if exists {
			toAssemble[name] = man
		}
	}

	err = a.Done(toAssemble)

//	dir, err := ioutil.TempDir("", "chunk-")
//	if err != nil {
//		return
//	}
//	defer os.RemoveAll(dir)
//
//	wg := sync.WaitGroup{}
//	var errs []error
//	for _, chunk := range chunkMap {
//		wg.Add(1)
//		go func(ch proto.ChunkInfo) {
//			defer wg.Done()
//			if err1 := t.FetchChunk(ch, dir); err1 != nil {
//				errs = append(errs, err1)
//				return
//			}
//		}(chunk)
//	}
// 	wg.Wait()
//	if len(errs) > 0 {
//		err = fmt.Errorf("errors while fetching chunks %s", errs)
//		return
//	}
//
//	// collect and relocate all
//	wg = sync.WaitGroup{}
//	for n, m := range blobs {
//		wg.Add(1)
//		go func(name string, man manifest.Manifest) {
//			defer wg.Done()
//			if err1 := t.collectBLOB(name, man, dir); err1 != nil {
//				errs = append(errs, err1)
//			}
//		}(n, m)
//	}
//	wg.Wait()
//	if len(errs) > 0 {
//		err = fmt.Errorf("errors while collecting blobs %s", errs)
//		return
//	}
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