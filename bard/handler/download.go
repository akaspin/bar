package handler
import (
"github.com/akaspin/bar/bard/storage"
"net/http"
	"fmt"
	"io"
"github.com/tamtam-im/logx"
)

type DownloadHandler struct {
	Storage *storage.StoragePool
	Prefix string
}

func (h *DownloadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	var id string
	if _, err = fmt.Sscanf(r.URL.Path, h.Prefix + "%s", &id); err != nil {
		logx.Error(err)
		w.WriteHeader(500)
		return
	}
	s, err := h.Storage.Take()
	if err != nil {
		logx.Error(err)
		w.WriteHeader(500)
		return
	}
	defer h.Storage.Release(s)

	rc, err := s.ReadBLOB(id)
	if err != nil {
		logx.Error(err)
		w.WriteHeader(404)
		return
	}
	defer rc.Close()

	n, err := io.Copy(w, rc)
	if err != nil {
		logx.Error(err)
		w.WriteHeader(500)
	}
	logx.Debugf("download %s: %d bytes sent", id, n)
}
