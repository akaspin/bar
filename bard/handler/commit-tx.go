package handler
import (
	"github.com/akaspin/bar/bard/storage"
	"net/http"
	"encoding/json"
	"github.com/akaspin/bar/proto"
	"github.com/tamtam-im/logx"
)

// Declare commit takes
type DeclareCommitTxHandler struct {
	Storage *storage.StoragePool
}

func (h *DeclareCommitTxHandler) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	var err error

	s, err := h.Storage.Take()
	if err != nil {
		logx.Error(err)
		w.WriteHeader(500)
		return
	}
	defer h.Storage.Release(s)

	req := proto.DeclareUploadTxRequest{}
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		logx.Error(err)
		w.WriteHeader(500)
		return
	}

	res := proto.DeclareUploadTxResponse{}
	for _, id := range req.IDs {
		if ok, _ := s.IsExists(id); !ok {
			res.MissingIDs = append(res.MissingIDs, id)
		}
	}
	if err = json.NewEncoder(w).Encode(&res); err != nil {
		logx.Error(err)
		w.WriteHeader(500)
		return
	}
}
