package transport
import (
	"github.com/akaspin/bar/proto"
	"time"
	"github.com/akaspin/bar/proto/manifest"
	"github.com/akaspin/bar/barc/model"
	"os"
	"sync"
	"github.com/tamtam-im/logx"
	"fmt"
	"path/filepath"
)

// Common transport with pooled connections
type Transport struct {
	WD string

	// base endpoint. http://example.com/v1
	DefaultEndpoint string

	rpcPool *RPCPool
}

// New RPC pool with default endpoint
func NewTransport(wd string, endpoint string, n int) (res *Transport) {
	res = &Transport{
		WD: wd,
		DefaultEndpoint: endpoint,
		rpcPool: NewRPCPool(n, time.Minute),
	}
	return
}

func (t *Transport) Close() {
	t.rpcPool.Close()
}

func (t *Transport) Ping() (res proto.Info, err error) {
	cli, err := t.rpcPool.Take(t.DefaultEndpoint)
	if err != nil {
		return
	}
	defer cli.Release()

	err = cli.Call("Service.Ping", &struct{}{}, &res)
	return
}

// Upload blobs
func (t *Transport) Upload(blobs model.Links) (err error) {
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

func (t *Transport) NewUpload(blobs model.Links) (toUpload model.Links, err error) {
	var req, res []manifest.Manifest

	idmap := blobs.IDMap()

	for _, name := range idmap {
		req = append(req, blobs[name])
	}

	cli, err := t.rpcPool.Take(t.DefaultEndpoint)
	if err != nil {
		return
	}
	defer cli.Release()

	if err = cli.Call("Service.NewUpload", &req, &res); err != nil {
		return
	}

	toUpload = model.Links{}
	for _, m := range res {
		toUpload[idmap[m.ID]] = m
	}
	return
}

func (t *Transport) UploadBlob(name string, info manifest.Manifest) (err error) {
	wg := sync.WaitGroup{}
	var errs []error

	logx.Debugf("uploading %s %s (%d chunks)", name, info.ID, len(info.Chunks))
	for _, chunk := range info.Chunks {
		wg.Add(1)
		go func(ci manifest.Chunk) {
			defer wg.Done()
			var err1 error
			if err1 = t.UploadChunk(name, info.ID, ci); err1 != nil {
				errs = append(errs, err1)
				return
			}

		}(chunk)
	}
	wg.Wait()
	if len(errs) > 0 {
		err = fmt.Errorf("errors while upload %s %s %s", name, info.ID, errs)
		return
	}
	return
}

func (t *Transport) FinishUpload(id string) (err error) {
	cli, err := t.rpcPool.Take(t.DefaultEndpoint)
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

func (t *Transport) UploadChunk(name string, blobID string, chunkInfo manifest.Chunk) (err error) {
	cli, err := t.rpcPool.Take(t.DefaultEndpoint)
	if err != nil {
		return
	}
	defer cli.Release()

	// read chunk
	f, err := os.Open(filepath.Join(t.WD, name))
	if err != nil {
		return
	}
	defer f.Close()

	buf := make([]byte, chunkInfo.Size)
	if _, err = f.ReadAt(buf, chunkInfo.Offset); err != nil {
		return
	}

	var res struct{}
	err = cli.Call("Service.UploadChunk", &proto.BLOBChunk{
		blobID, chunkInfo.ID, chunkInfo.Size, buf,
	}, &res)

	logx.Debugf("chunk %s %s:%s %d uploaded", name,
		blobID, chunkInfo.ID, chunkInfo.Size)
	return
}

