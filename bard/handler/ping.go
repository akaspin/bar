package handler
import (
	"net/http"
	"encoding/json"
	"github.com/tamtam-im/logx"
)

//
type PingHandler struct {
	ChunkSize int64
}

func (h *PingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := struct{
		ChunkSize int64
	}{h.ChunkSize}
	resp, err := json.Marshal(&data)
	if err != nil {
		logx.Error(err)
		w.WriteHeader(500)
	}
	_, err = w.Write(resp)
	logx.OnError(err)
	logx.Debugf("pong %v", data)
}
