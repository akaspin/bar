package handler
import (
	"net/http"
	"github.com/akaspin/bar/bard/storage"
	"fmt"
)

// Just accepts simple uploads
type SimpleUploadHandler struct {
	Storage *storage.StoragePool
}

func (h *SimpleUploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	s, err := h.Storage.Take()
	if err != nil {
		w.WriteHeader(500)
		return
	}
	defer h.Storage.Release(s)

	var id string
	if _, err = fmt.Sscanf(r.URL.Path, "/v1/blob/upload/%s", &id); err != nil {
		w.WriteHeader(500)
		return
	}

	err = s.StoreBLOB(id, r.ContentLength, r.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
}
