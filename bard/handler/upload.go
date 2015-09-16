package handler
import (
	"net/http"
	"github.com/akaspin/bar/bard/storage"
)

// Just accepts simple uploads
type SimpleUploadHandler struct {
	storagePool *storage.StoragePool
}

func (h *SimpleUploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request)  {

}
