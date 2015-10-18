package handler
import (
	"github.com/akaspin/bar/proto"
	"net/http"
	h_template "html/template"
	"github.com/tamtam-im/logx"
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
	<a href="{{.Info.HTTPEndpoint}}/win/bar-export.bat"><code>bar-export.bat</code></a>
	and save in root of the working tree.</p>
<p>This script is no-brain solution to export bar specs. It automatically
download <code>bar.exe</code> if it is not found and upload BLOBs and spec
to bard</p>
<p>You can also download <a href="{{.Info.HTTPEndpoint}}/win/bar.exe"><code>bar.exe</code></a> and save
	it beside <code>bar-export.bat</code> or somewhere in PATH.</p>

<h2>Git users</h2>
<p>To install bar into git repository use <code>bar git-inir</code></p>
<pre>
$ bar -endpoint={{.Info.JoinRPCEndpoints}} git-init
</pre>
</body>
</html>
`

type FrontHandler struct {
	Info *proto.ServerInfo
}

func (h *FrontHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	front_template, err := h_template.New("front").Parse(front_tpl)
	if err != nil {
		logx.Error(err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err = front_template.Execute(w, map[string]interface{}{
		"Info": h.Info,
	}); err != nil {
		logx.Error(err)
		w.WriteHeader(500)
		return
	}
}