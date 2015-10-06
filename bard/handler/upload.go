package handler
import (
	"net/http"
	"github.com/akaspin/bar/bard/storage"
	"fmt"
	"github.com/tamtam-im/logx"
)

// Just accepts simple uploads
type SimpleUploadHandler struct {
	Storage *storage.StoragePool
	Prefix string
}

func (h *SimpleUploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	s, err := h.Storage.Take()
	if err != nil {
		logx.Error(err)
		w.WriteHeader(500)
		return
	}
	defer h.Storage.Release(s)

	var id string
	if _, err = fmt.Sscanf(r.URL.Path, h.Prefix + "%s", &id); err != nil {
		logx.Error(err)
		w.WriteHeader(500)
		return
	}

	var size int64
	if _, err = fmt.Sscanf(r.Header.Get("blob-size"), "%d", &size); err != nil {
		logx.Error(err)
		w.WriteHeader(500)
		return
	}
	err = s.WriteBLOB(id, size, r.Body)
	if err != nil {
		logx.Error(err)
		w.WriteHeader(500)
		return
	}
	logx.Infof("stored %s, %d bytes", id, size)
	w.WriteHeader(200)
}
