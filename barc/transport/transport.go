package transport
import (
	"github.com/akaspin/bar/proto"
	"time"
	"github.com/akaspin/bar/proto/manifest"
	"github.com/akaspin/bar/barc/model"
)

// Common transport with pooled connections
type Transport struct {
	// base endpoint. http://example.com/v1
	DefaultEndpoint string

	rpcPool *RPCPool
}

// New RPC pool with default endpoint
func NewTransport(endpoint string, n int) (res *Transport) {
	res = &Transport{DefaultEndpoint: endpoint, rpcPool: NewRPCPool(n, time.Minute)}
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
	cli, err := t.rpcPool.Take(t.DefaultEndpoint)
	if err != nil {
		return
	}
	defer cli.Release()

	// Make idmap to deduplicate uploads
	idmap := map[string]string{}
	for name, m := range blobs {
		idmap[m.ID] = name
	}

	// Declare blobs on bard
	var declareReq []manifest.Manifest
	for id, _ := range idmap {
		declareReq = append(declareReq, *blobs[idmap[id]])
	}



	return
}

func (t *Transport) NewUpload(blobs model.Links) (toUpload model.Links, err error) {
	var req, res []manifest.Manifest

	idmap := blobs.IDMap()

	for _, name := range idmap {
		req = append(req, *blobs[name])
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
		toUpload[idmap[m.ID]] = &m
	}

	return
}

