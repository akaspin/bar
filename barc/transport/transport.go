package transport
import (
	"github.com/akaspin/bar/proto"
	"time"
	"github.com/akaspin/bar/proto/manifest"
	"os"
	"sync"
	"github.com/tamtam-im/logx"
	"fmt"
	"path/filepath"
	"io/ioutil"
	"github.com/akaspin/go-contentaddressable"
	"golang.org/x/crypto/sha3"
	"io"
	"github.com/akaspin/bar/barc/lists"
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
		rpcPool: NewRPCPool(n, time.Minute * 30),
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

func (t *Transport) NewUpload(blobs lists.Links) (toUpload lists.Links, err error) {
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

	toUpload = lists.Links{}
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
	err = cli.Call("Service.UploadChunk", &proto.ChunkData{
		proto.ChunkInfo{blobID, chunkInfo}, buf,
	}, &res)

	logx.Tracef("chunk %s %s:%s %d uploaded", name,
		blobID, chunkInfo.ID, chunkInfo.Size)
	return
}

func (t *Transport) Download(blobs lists.Links) (err error) {

	fetch, err := t.GetFetch(blobs.IDMap().IDs())
	logx.Debug("fetching blobs %s", blobs.IDMap())

	// little funny, but all chunks are equal - flatten them
	chunkMap := map[string]proto.ChunkInfo{}
	for _, mt := range fetch {
		for _, ch := range mt.Chunks {
			chunkMap[ch.ID] = proto.ChunkInfo{mt.ID, ch}
		}
	}

	// Fetch all chunks to temporary dir
	dir, err := ioutil.TempDir("", "chunk-")
	if err != nil {
		return
	}
	defer os.RemoveAll(dir)

	wg := sync.WaitGroup{}
	var errs []error
	for _, chunk := range chunkMap {
		wg.Add(1)
		go func(ch proto.ChunkInfo) {
			defer wg.Done()
			if err1 := t.FetchChunk(ch, dir); err1 != nil {
				errs = append(errs, err1)
				return
			}
		}(chunk)
	}
 	wg.Wait()
	if len(errs) > 0 {
		err = fmt.Errorf("errors while fetching chunks %s", errs)
		return
	}

	// collect and relocate all
	wg = sync.WaitGroup{}
	for n, m := range blobs {
		wg.Add(1)
		go func(name string, man manifest.Manifest) {
			defer wg.Done()
			if err1 := t.collectBLOB(name, man, dir); err1 != nil {
				errs = append(errs, err1)
			}
		}(n, m)
	}
	wg.Wait()
	if len(errs) > 0 {
		err = fmt.Errorf("errors while collecting blobs %s", errs)
		return
	}
	return
}

func (t *Transport) collectBLOB(name string, man manifest.Manifest, dir string) (err error) {
	logx.Debugf("assembling blob %s %s", name, man.ID)

	blobdir, err := ioutil.TempDir("", "")
	if err != nil {
		return
	}
	defer os.RemoveAll(blobdir)

	caOpts := contentaddressable.DefaultOptions()
	caOpts.Hasher = sha3.New256()

	w, err := contentaddressable.NewFileWithOptions(
		filepath.Join(blobdir, man.ID), caOpts)
	if err != nil {
		return
	}
	defer w.Close()

	for _, chunk := range man.Chunks {
		written, err := t.writeChunkTo(filepath.Join(dir, chunk.ID), w)
		if err != nil {
			return err
		}
		if written != chunk.Size {
			return fmt.Errorf("bad size for chunk %s %s %s %d != %d",
				name, man.ID, chunk.ID, chunk.Size, written)
		}
	}
	if err = w.Accept(); err != nil {
		return
	}
	os.Remove(filepath.Join(t.WD, name))
	if err = os.Rename(filepath.Join(blobdir, man.ID), filepath.Join(t.WD, name)); err != nil {
		return
	}
	logx.Debugf("done assemble blob %s %s", name, man.ID)

	return
}

func (t *Transport) writeChunkTo(src string, dst io.Writer) (n int64, err error) {
	r, err := os.Open(src)
	if err != nil {
		return
	}
	defer r.Close()
	n, err = io.Copy(dst, r)
	return
}

func (t *Transport) GetFetch(ids []string) (res []manifest.Manifest, err error) {
	cli, err := t.rpcPool.Take(t.DefaultEndpoint)
	if err != nil {
		return
	}
	defer cli.Release()
	err = cli.Call("Service.GetFetch", &ids, &res)
	return
}

func (t *Transport) FetchChunk(info proto.ChunkInfo, dir string) (err error) {
	logx.Debugf("fetching chunk %s", info.ID)

	caOpts := contentaddressable.DefaultOptions()
	caOpts.Hasher = sha3.New256()

	w, err := contentaddressable.NewFileWithOptions(
		filepath.Join(dir, info.ID), caOpts)
	if err != nil {
		return
	}
	defer w.Close()

	cli, err := t.rpcPool.Take(t.DefaultEndpoint)
	if err != nil {
		return
	}
	defer cli.Release()

	var res proto.ChunkData
	if err = cli.Call("Service.FetchChunk", &info, &res); err != nil {
		return
	}
	if _, err = w.Write(res.Data); err != nil {
		return
	}
	if err = w.Accept(); err != nil {
		logx.Error(err)
		return
	}
	logx.Debugf("done fetch chunk %s", info.ID)
	return
}

func (t *Transport) Check(ids []string) (res []string, err error) {
	cli, err := t.rpcPool.Take(t.DefaultEndpoint)
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
	cli, err := t.rpcPool.Take(t.DefaultEndpoint)
	if err != nil {
		return
	}
	defer cli.Release()

	var res struct{}
	err = cli.Call("Service.UploadSpec", &spec, &res)

	return
}

func (t *Transport) GetSpec(id string) (res lists.Links, err error) {
	cli, err := t.rpcPool.Take(t.DefaultEndpoint)
	if err != nil {
		return
	}
	defer cli.Release()
	err = cli.Call("Service.GetSpec", &id, &res)
	return
}