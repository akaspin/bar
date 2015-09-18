package handler
import (
"github.com/akaspin/bar/bard/storage"
"net/http"
	"fmt"
)

type CheckHandler struct {
	Storage *storage.StoragePool
	Prefix string
}

func (h *CheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	var id string
	var err error

	if _, err = fmt.Sscanf(r.URL.Path, h.Prefix + "%s", &id); err != nil {
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
}
