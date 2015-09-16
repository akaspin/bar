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

	var size int64
	if _, err = fmt.Sscanf(r.Header.Get("BLOB-Size"), "%d", &size); err != nil {
		w.WriteHeader(500)
		return
	}
	err = s.StoreBLOB(r.Header.Get("BLOB-ID"), size, r.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
}
