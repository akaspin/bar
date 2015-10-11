package handler
import (
	"net/http"
	"strings"
	"github.com/tamtam-im/logx"
	"github.com/akaspin/bar/proto"
	"text/template"

	"bytes"
	"os"
	"io"
)

const export_bat_tpl = `
@echo off

@WHERE barc
IF %ERRORLEVEL% NEQ 0 (
	ECHO barc is not found. downloading...
	powershell -command "$clnt = new-object System.Net.WebClient; $clnt.DownloadFile(\"{{.HTTPEndpoint}}/win/barc.exe\", \"barc.exe\")"
)

barc -log-level=DEBUG up -endpoint={{.Endpoint}} -chunk={{.ChunkSize}} -pool={{.PoolSize}} !bar*.bat !bar-spec*.json !barc.exe !desktop.ini
for /f %%i in ('barc -log-level=DEBUG spec-export -endpoint={{.Endpoint}} -chunk={{.ChunkSize}} -pool={{.PoolSize}} -upload -cc !bar*.bat !bar-spec*.json !barc.exe !desktop.ini') do set VAR=%%i

start {{.HTTPEndpoint}}/spec/%VAR%

echo press any key...
pause >nul
`

const import_bat_tpl = `
@echo off

@WHERE barc
IF %ERRORLEVEL% NEQ 0 (
	ECHO barc is not found. downloading...
	powershell -command "$clnt = new-object System.Net.WebClient; $clnt.DownloadFile(\"{{.HTTPEndpoint}}/win/barc.exe\", \"barc.exe\")"
)

barc -log-level=DEBUG spec-import -endpoint={{.Endpoint}} -chunk={{.ChunkSize}} -pool={{.PoolSize}} {{.ID}}
barc -log-level=DEBUG down -endpoint={{.Endpoint}} -chunk={{.ChunkSize}} -pool={{.PoolSize}} {{.ID}}

echo press any key...
pause >nul
`

type ExportBatHandler struct {
	Info *proto.Info
}

func (h *ExportBatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bat_template, err := template.New("bat_template").Parse(export_bat_tpl)
	if err != nil {
		return
	}
	buf := new(bytes.Buffer)

	if err = bat_template.Execute(buf, h.Info); err != nil {
		return
	}
	body := string(buf.Bytes())
	crlf := strings.Replace(body, "\n", "\r\n", -1)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(crlf))
}

type ImportBatHandler struct {
	Info *proto.Info
	BarExe string
}

func (h *ImportBatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/win/bar-import/")[:64]

	w.Write([]byte(id))
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


