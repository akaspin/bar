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

@WHERE bar
IF %ERRORLEVEL% NEQ 0 (
	ECHO bar is not found. downloading...
	powershell -command "$clnt = new-object System.Net.WebClient; $clnt.DownloadFile(\"{{.HTTPEndpoint}}/win/bar.exe\", \"bar.exe\")"
)

bar -log-level=DEBUG up -http={{.HTTPEndpoint}} -rpc={{.JoinRPCEndpoints}} -chunk={{.ChunkSize}} -pool={{.PoolSize}} !bar*.bat !bar-spec*.json !bar.exe !desktop.ini
for /f %%i in ('bar -log-level=DEBUG spec-export -http={{.HTTPEndpoint}} -rpc={{.JoinRPCEndpoints}} -chunk={{.ChunkSize}} -pool={{.PoolSize}} -upload -cc !bar*.bat !bar-spec*.json !bar.exe !desktop.ini') do set VAR=%%i

start {{.HTTPEndpoint}}/spec/%VAR%

echo press any key...
pause >nul
`

const import_bat_tpl = `
@echo off

@WHERE bar
IF %ERRORLEVEL% NEQ 0 (
	ECHO bar is not found. downloading...
	powershell -command "$clnt = new-object System.Net.WebClient; $clnt.DownloadFile(\"{{.Info.HTTPEndpoint}}/win/bar.exe\", \"bar.exe\")"
)

for /f %%i in ('bar -log-level=DEBUG spec-import -http={{.Info.HTTPEndpoint}} -rpc={{.Info.JoinRPCEndpoints}} -chunk={{.Info.ChunkSize}} -pool={{.Info.PoolSize}} {{.ID}}') do set VAR=%%i
bar -log-level=DEBUG down -http={{.Info.HTTPEndpoint}} -rpc={{.Info.JoinRPCEndpoints}} -chunk={{.Info.ChunkSize}} -pool={{.Info.PoolSize}} %VAR%

echo press any key...
pause >nul
`

type ExportBatHandler struct {
	Info *proto.ServerInfo
}

func (h *ExportBatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bat_template, err := template.New("bat_template").Parse(export_bat_tpl)
	if err != nil {
		logx.Error(err)
		return
	}
	buf := new(bytes.Buffer)

	if err = bat_template.Execute(buf, h.Info); err != nil {
		logx.Error(err)
		return
	}
	body := string(buf.Bytes())
	crlf := strings.Replace(body, "\n", "\r\n", -1)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(crlf))
}

type ImportBatHandler struct {
	Info *proto.ServerInfo
}

func (h *ImportBatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/win/bar-import/")[:64]

	bat_template, err := template.New("bat_template").Parse(import_bat_tpl)
	if err != nil {
		return
	}
	buf := new(bytes.Buffer)

	if err = bat_template.Execute(buf, map[string]interface{}{
		"Info": h.Info,
		"ID": id,
	}); err != nil {
		return
	}
	body := string(buf.Bytes())
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


