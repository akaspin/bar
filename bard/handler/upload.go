package handler
import (
	"net/http"
	"github.com/akaspin/bar/bard/storage"
	"fmt"
)

// Just accepts simple uploads
type SimpleUploadHandler struct {
	Storage *storage.StoragePool
	Prefix string
}

func (h *SimpleUploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	s, err := h.Storage.Take()
	if err != nil {
		w.WriteHeader(500)
		return
	}
	defer h.Storage.Release(s)

	var id string
	if _, err = fmt.Sscanf(r.URL.Path, h.Prefix + "%s", &id); err != nil {
		w.WriteHeader(500)
		return
	}

	var size int64
	if _, err = fmt.Sscanf(r.Header.Get("blob-size"), "%d", &size); err != nil {
		w.WriteHeader(500)
		return
	}
	err = s.StoreBLOB(id, size, r.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
}
