package handler
import (
"github.com/akaspin/bar/bard/storage"
"net/http"
)

type SpecUploadHandler struct {
	Storage *storage.StoragePool
}

func (h *SpecUploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}