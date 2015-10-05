package transport
import (
	"github.com/akaspin/bar/proto/manifest"
	"fmt"
	"net/url"
	"path"
	"net/http"
	"os"
	"io/ioutil"
	"encoding/json"
	"bytes"
	"io"
	"github.com/akaspin/bar/proto"
	"github.com/tamtam-im/logx"
)


type Transport struct {
	// bard endpoint. http://example.com/v1
	Endpoint *url.URL
}

func (t *Transport) Ping() (res proto.Info, err error) {
	api := t.apiURL("ping").String()
	logx.Debugf("sending ping to %s", api)
	resp, err := http.Get(api)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("bad pong: %d %s", resp.StatusCode, resp.Status)
	}
	res = proto.Info{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	logx.Debugf("pong received from %s %v", api, res)
	return
}

// Declare commit transaction and get existent ids.
// This similar to Transport.Check but declares new git commit. See DoneCommit
func (t *Transport) DeclareCommitTx(txID string, ids []string) (res []string, err error) {
	buf := new(bytes.Buffer)

	if err = json.NewEncoder(buf).Encode(proto.DeclareUploadTxRequest{txID, ids}); err != nil {
		return
	}

	resp, err := http.Post(t.apiURL("/tx/commit/declare").String(),
		"application/octet-stream", bytes.NewReader(buf.Bytes()))
	if err != nil {
		return
	}

	txResp := proto.DeclareUploadTxResponse{}
	if err = json.NewDecoder(resp.Body).Decode(&txResp); err != nil {
		return
	}
	defer resp.Body.Close()
	res = txResp.MissingIDs
	return
}

// Done commit transaction
func (t *Transport) DoneCommitTx(txID, commitID string) (err error) {
	// TODO
	return
}

func (t *Transport) AbortCommitTx(txID string) (err error) {
	// TODO
	return
}

// Push BLOB regardless of index state. Filename MUST be relative
func (t *Transport) Push(filename string, manifest *manifest.Manifest) (err error) {
	r, err := os.Open(filename)
	if err != nil {
		return
	}
	defer r.Close()

	req, err := http.NewRequest("POST",
		t.apiURL(fmt.Sprintf("/blob/upload/%s", manifest.ID)).String(),
		r,
	)
	req.Close = true
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("blob-size", fmt.Sprintf("%d", manifest.Size))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	return
}

func (t *Transport) Check(ids []string) (res []string, err error) {
	buf, err := json.Marshal(ids)

	resp, err := http.Post(t.apiURL("/blob/check").String(),
		"application/octet-stream", bytes.NewReader(buf))
	if err != nil {
		return
	}

	bodyBuf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = json.Unmarshal(bodyBuf, &res)
	return
}

// Get BLOB from bard
func (t *Transport) GetBLOB(id string, size int64, w io.Writer) (err error) {
	resp, err := http.Get(t.apiURL(
		fmt.Sprintf("/blob/download/%s", id)).String())
	if err != nil {
		return
	}
	defer resp.Body.Close()

	n, err := io.CopyN(w, resp.Body, size)
	if err != nil {
		return
	}
	if n != size {
		err = fmt.Errorf("bad download size for %s: expect %d, got %d",
			id, size, n)
	}
	return
}

func (t *Transport) Info(id string) (err error) {
	req := &http.Request{
		Method: "GET",
		URL: t.apiURL(fmt.Sprintf("/blob/info/%s", id)),
	}
	_, err = http.DefaultClient.Do(req)
	return
}

func (t *Transport) apiURL(meth string) (res *url.URL) {
	res = new(url.URL)
	*res = *t.Endpoint
	res.Path = path.Join(res.Path, meth)
	return
}
