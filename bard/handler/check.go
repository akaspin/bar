package handler
import (
	"github.com/akaspin/bar/bard/storage"
	"net/http"
	"io/ioutil"
	"strings"
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

	res := []string{}
	for _, id := range strings.Split(string(buf), "\n") {
		if id == "" {
			continue
		}
		if ok, _ := s.IsExists(id); ok {
			res = append(res, id)
		}
	}

	w.Write([]byte(strings.Join(res, "\n")))
}
