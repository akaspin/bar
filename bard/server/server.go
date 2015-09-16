package server
import (
	"github.com/akaspin/bar/bard/storage"
	"net/http"
	"github.com/akaspin/bar/bard/handler"
)

func Serve(addr string, storagePool *storage.StoragePool) (err error) {
	mux := http.NewServeMux()

	mux.Handle("/v1/blob/upload/", &handler.SimpleUploadHandler{storagePool})
	mux.Handle("/v1/blob/info/", &handler.InfoHandler{storagePool})

	s := &http.Server{Addr:addr, Handler: mux}
	err = s.ListenAndServe()
	return
}
