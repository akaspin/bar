package front

import (
	"bytes"
	"github.com/akaspin/bar/proto"
	"github.com/akaspin/bar/server/storage"
	"github.com/tamtam-im/logx"
	"golang.org/x/net/context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Handlers struct {
	ctx     context.Context
	options *Options
	storage.Storage

	*HttpTpl
}

func NewHandlers(ctx context.Context, options *Options, s storage.Storage) (res *Handlers, err error) {
	res = &Handlers{
		ctx:     ctx,
		options: options,
		Storage: s,
	}
	res.HttpTpl, err = NewHttpTpls()
	return
}

// http://_:_/
func (h *Handlers) HandleFront(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	h.handleTpl(w, "front", map[string]interface{}{
		"Info": h.options.Info,
	})
}

func (h *Handlers) HandleSpec(w http.ResponseWriter, r *http.Request) {
	id := proto.ID(strings.TrimPrefix(r.URL.Path, "/v1/spec/"))

	logx.Debugf("serving spec %s", id)

	ok, err := h.Storage.IsSpecExists(id)
	if err != nil {
		logx.Error(err)
		w.WriteHeader(500)
		return
	}
	if !ok {
		logx.Errorf("bad spec id %s", id)
		w.WriteHeader(404)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	h.handleTpl(w, "spec", map[string]interface{}{
		"Info":    h.options.Info,
		"ID":      id,
		"ShortID": id[:12],
	})
}

func (h *Handlers) HandleExportBat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	h.handleWinTpl(w, "bar-export.bat", map[string]interface{}{
		"Info": h.options.Info,
	})
}

func (h *Handlers) HandleImportBat(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/win/bar-import/")[:64]
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	h.handleWinTpl(w, "bar-import.bat", map[string]interface{}{
		"Info": h.options.Info,
		"ID":   id,
	})
}

func (h *Handlers) HandleBarExe(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open(filepath.Join(h.options.BinDir, "windows"))
	if err != nil {
		logx.Error(err)
		w.WriteHeader(500)
		return
	}
	defer f.Close()

	_, err = io.Copy(w, f)
	if err != nil {
		logx.Error(err)
		w.WriteHeader(500)
		return
	}
}

func (h *Handlers) handleWinTpl(w http.ResponseWriter, name string, data map[string]interface{}) (err error) {
	buf := new(bytes.Buffer)

	if err = h.HttpTpl.tpls[name](data, w); err != nil {
		logx.Error(err)
		w.WriteHeader(500)
		return
	}
	body := string(buf.Bytes())
	crlf := strings.Replace(body, "\n", "\r\n", -1)
	w.Write([]byte(crlf))
	return
}

func (h *Handlers) handleTpl(w http.ResponseWriter, name string, data map[string]interface{}) (err error) {
	if err = h.HttpTpl.tpls[name](data, w); err != nil {
		logx.Error(err)
		w.WriteHeader(500)
	}
	return
}
