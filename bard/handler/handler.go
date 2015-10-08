package handler
import (
	"net/http"
	"github.com/akaspin/bar/bard/storage"
	"strings"
	"github.com/tamtam-im/logx"
	"encoding/json"
)


type FrontHandler struct {

}

func (h *FrontHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("works"))
}

//
type BatHandler struct {

}

func (h *BatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("works"))
}


type SpecHandler struct {
	Storage *storage.StoragePool
}

func (h *SpecHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	store, err := h.Storage.Take()
	if err != nil {
		return
	}
	defer h.Storage.Release(store)

	id := strings.TrimPrefix(r.URL.Path, "/v1/spec/")

	logx.Debugf("serving spec %s", id)

	spec, err := store.ReadSpec(id)
	if err != nil {
		logx.Errorf("bad spec id %s", id)
		w.WriteHeader(404)
		return
	}

	w.Header().Set("Content-Type", "application/javascript")

	if err = json.NewEncoder(w).Encode(&spec); err != nil {
		logx.Error(err)
		w.WriteHeader(500)
		return
	}
}
