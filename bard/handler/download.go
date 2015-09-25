package handler
import (
"github.com/akaspin/bar/bard/storage"
"net/http"
	"fmt"
	"io"
)

type DownloadHandler struct {
	Storage *storage.StoragePool
	Prefix string
}

func (h *DownloadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	var id string
	if _, err = fmt.Sscanf(r.URL.Path, h.Prefix + "%s", &id); err != nil {
		w.WriteHeader(500)
		return
	}
	s, err := h.Storage.Take()
	if err != nil {
		w.WriteHeader(500)
		return
	}
	defer h.Storage.Release(s)

	rc, err := s.ReadBLOB(id)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	defer rc.Close()

	_, err = io.Copy(w, rc)
	if err != nil {
		w.WriteHeader(500)
	}
}
