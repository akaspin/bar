package handler
import (
	"net/http"
	"github.com/akaspin/bar/bard/storage"
	"strings"
	"github.com/tamtam-im/logx"
	"encoding/json"
	"github.com/akaspin/bar/proto"
	"text/template"
	h_template "html/template"
	"bytes"
	"os"
	"io"
)

const front_tpl = `
<!DOCTYPE html>
<html>
<head>
<title>BAR</title>
</head>
<body>
<h1>Welcome to BAR!</h1>
<p>BAR is simple BLOB vendoring system.
	Visit <a href="https://github.com/akaspin/bar">github repo</a> for details.</p>

<h2>Windows users</h2>
<p>To export specs without git download
	<a href="{{.Info.HTTPEndpoint}}/bar-export.bat"><code>bar-export.bat</code></a> or save next
	code as bat-file in root of the working tree.</p>
<pre>
{{.BatBody}}
</pre>
<p>This script is no-brain solution to export bar specs. It automatically
download <code>barc.exe</code> if it is not found and upload BLOBs and spec
to bard</p>
<p>You can also download <a href="{{.Info.HTTPEndpoint}}/barc.exe"><code>barc.exe</code></a> and save
	it beside <code>bar-export.bat</code> or somewhere in PATH.</p>
</body>
</html>
`

const bat_tpl = `
@echo off

@WHERE barc
IF %ERRORLEVEL% NEQ 0 (
	ECHO barc is not found. downloading...
	powershell -command "$clnt = new-object System.Net.WebClient; $clnt.DownloadFile(\"http://{{.HTTPEndpoint}}:3000/v1/barc.exe\", \"barc.exe\")"
)

barc -log-level=DEBUG up -endpoint={{.Endpoint}} -chunk={{.ChunkSize}} -pool={{.PoolSize}} !bar*.bat !bar-spec*.json !barc.exe !desktop.ini
for /f %%i in ('barc -log-level=DEBUG spec-export -endpoint={{.Endpoint}} -chunk={{.ChunkSize}} -pool={{.PoolSize}} -upload -cc !bar*.bat !bar-spec*.json !barc.exe !desktop.ini') do set VAR=%%i

start {{.HTTPEndpoint}}/spec/%VAR%

echo press any key...
pause >nul
`

func batBody(info *proto.Info) (res string, err error) {
	bat_template, err := template.New("bat_template").Parse(bat_tpl)
	if err != nil {
		return
	}
	buf := new(bytes.Buffer)

	if err = bat_template.Execute(buf, info); err != nil {
		return
	}
	res = string(buf.Bytes())
	return
}

type FrontHandler struct {
	Info *proto.Info
	BarExe string
}

func (h *FrontHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	front_template, err := h_template.New("front").Parse(front_tpl)
	if err != nil {
		logx.Error(err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	body, err := batBody(h.Info)
	if err != nil {
		logx.Error(err)
		w.WriteHeader(500)
		return
	}
	if err = front_template.Execute(w, map[string]interface{}{
		"Info": h.Info,
		"BatBody": body,
		"BarExe": h.BarExe,
	}); err != nil {
		logx.Error(err)
		w.WriteHeader(500)
		return
	}
}

//
type BatHandler struct {
	Info *proto.Info
}

func (h *BatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := batBody(h.Info)
	if err != nil {
		logx.Error(err)
		w.WriteHeader(500)
		return
	}
	crlf := strings.Replace(body, "\n", "\r\n", -1)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(crlf))
}

type ExeHandler struct {
	Exe string
}

func (h *ExeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open(h.Exe)
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


type SpecHandler struct {
	Storage *storage.StoragePool
}

func (h *SpecHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	store, err := h.Storage.Take()
	if err != nil {
		return
	}
	defer h.Storage.Release(store)

	id := strings.TrimPrefix(r.URL.Path, "/v1/spec/")

	logx.Debugf("serving spec %s", id)

	spec, err := store.ReadSpec(id)
	if err != nil {
		logx.Errorf("bad spec id %s", id)
		w.WriteHeader(404)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	if err = json.NewEncoder(w).Encode(&spec); err != nil {
		logx.Error(err)
		w.WriteHeader(500)
		return
	}
}
