package handler
import (
	"net/http"
	"encoding/json"
	"github.com/tamtam-im/logx"
	"github.com/akaspin/bar/proto"
)

//
type PingHandler struct {
	ChunkSize int64
	ClientConns int
}

func (h *PingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := proto.Info{h.ChunkSize, h.ClientConns}
	resp, err := json.Marshal(&data)
	if err != nil {
		logx.Error(err)
		w.WriteHeader(500)
	}
	_, err = w.Write(resp)
	logx.OnError(err)
	logx.Debugf("pong %v", data)
}
