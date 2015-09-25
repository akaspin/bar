package server
import (
	"github.com/akaspin/bar/bard/storage"
	"net/http"
	"github.com/akaspin/bar/bard/handler"
)

func Serve(addr string, storagePool *storage.StoragePool) (err error) {
	mux := http.NewServeMux()

	mux.Handle("/v1/blob/upload/", &handler.SimpleUploadHandler{
		storagePool, "/v1/blob/upload/"})
	mux.Handle("/v1/blob/check", &handler.CheckHandler{
		storagePool})
	mux.Handle("/v1/tx/commit/declare/", &handler.DeclareCommitTxHandler{
		storagePool, "/v1/tx/commit/declare/"})
	mux.Handle("/v1/blob/download/", &handler.DownloadHandler{
		storagePool, "/v1/blob/download/"})

	s := &http.Server{Addr:addr, Handler: mux}
	err = s.ListenAndServe()
	return
}
