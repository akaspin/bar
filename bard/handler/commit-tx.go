package handler
import (
"github.com/akaspin/bar/bard/storage"
"net/http"
	"io/ioutil"
	"encoding/json"
)


type DeclareCommitTxHandler struct {
	Storage *storage.StoragePool
	Prefix string
}

func (h *DeclareCommitTxHandler) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	s, err := h.Storage.Take()
	if err != nil {
		w.WriteHeader(500)
		return
	}
	defer h.Storage.Release(s)

	var in, res []string
	err = json.Unmarshal(buf, &in)
	for _, id := range in {
		if ok, _ := s.IsExists(id); ok {
			res = append(res, id)
		}
	}
	out, err := json.Marshal(res)
	w.Write(out)
}
