package transport
import (
	"github.com/akaspin/bar/shadow"
	"fmt"
	"net/url"
	"path"
	"net/http"
	"encoding/hex"
	"os"
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

	req := &http.Request{
		Method: "POST",
		URL: t.apiURL("/blob/upload"),
		Header: http.Header{
			"BLOB-ID": {hex.EncodeToString(manifest.ID)},
			"BLOB-Size": {fmt.Sprintf("%d", manifest.Size)},
		},
		Body: r,
		ContentLength: manifest.Size,
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
