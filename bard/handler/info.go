package handler
import (
	"github.com/akaspin/bar/bard/storage"
	"net/http"
	"fmt"
)

// Just accepts simple uploads
type InfoHandler struct {
	Storage *storage.StoragePool
	Prefix string
}

func (h *InfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	var id string
	var err error

	if _, err = fmt.Sscanf(r.URL.Path, h.Prefix + "%s", &id); err != nil {
		w.WriteHeader(500)
		return
	}

//	var chunkSize int64
//	if r.Header.Get("blob-chunk-size") != "" {
//
//	}

	w.WriteHeader(200)
}
