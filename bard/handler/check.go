package handler
import (
	"github.com/akaspin/bar/bard/storage"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

// Check for existent blobs. Takes list of needed blobs and
// responses with list of existent blobs.
type CheckHandler struct {
	Storage *storage.StoragePool
}

func (h *CheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	s, err := h.Storage.Take()
	if err != nil {
		w.WriteHeader(500)
		return
	}
	defer h.Storage.Release(s)

	var in, res []string
	err = json.Unmarshal(buf, &in)
	for _, id := range in {
		if ok, _ := s.IsExists(id); ok {
			res = append(res, id)
		}
	}
	out, err := json.Marshal(res)
	w.Write(out)
}
