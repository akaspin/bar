package transport
import (
	"github.com/akaspin/bar/shadow"
	"fmt"
	"net/url"
	"path"
	"net/http"
	"os"
	"bytes"
	"io/ioutil"
	"strings"
	"encoding/hex"
)


type Transport struct {
	// bard endpoint. http://example.com/v1
	Endpoint *url.URL
}

// Push BLOB regardless of index state. Filename MUST be relative
func (t *Transport) Push(filename string, manifest *shadow.Shadow) (err error) {
	if manifest.IsFromShadow {
		err = fmt.Errorf("%x is shadow skip push")
	}

	r, err := os.Open(filename)
	if err != nil {
		return
	}
	defer r.Close()

	req, err := http.NewRequest("POST",
		t.apiURL(fmt.Sprintf("/blob/upload/%x", manifest.ID)).String(),
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

func (t *Transport) Check(ids [][]byte) (res [][]byte, err error) {
	buf := new(bytes.Buffer)

	for _, id := range ids {
		if _, err = buf.WriteString(fmt.Sprintf("%x\n", id)); err != nil {
			return
		}
	}

	resp, err := http.Post(t.apiURL("/blob/check").String(),
		"application/octet-stream", buf)
	if err != nil {
		return
	}
	bodyBuf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	for _, i := range strings.Split(string(bodyBuf), "\n") {
		var id []byte
		id, err = hex.DecodeString(i)
		if err != nil {
			return
		}
		res = append(res, id)
	}
	return
}

func (t *Transport) Info(id []byte) (err error) {
	req := &http.Request{
		Method: "GET",
		URL: t.apiURL(fmt.Sprintf("/blob/info/%x", id)),
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
