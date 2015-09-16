package handler
import (
"github.com/akaspin/bar/bard/storage"
"net/http"
	"log"
)

// Just accepts simple uploads
type InfoHandler struct {
	Storage *storage.StoragePool
}

func (h *InfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	log.Println(r.URL)
	w.WriteHeader(200)
}
