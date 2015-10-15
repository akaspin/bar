package handler
import (
"github.com/akaspin/bar/bard/storage"
"net/http"
"strings"
	"github.com/tamtam-im/logx"
	"encoding/json"
	h_template "html/template"
	"github.com/akaspin/bar/proto"
)

const spec_tpl string = `
<!DOCTYPE html>
<html>
<head>
<title>SPEC {{.ID}}</title>
</head>
<body>
<h1>:-)</h1>
<pre>
{{.ID}}
</pre>
<p>Send it to someone who knows what to do with it.</p>
<h2>Windows users</h2>
<p>To import spec download
	<a href="{{.Info.HTTPEndpoint}}/win/bar-import/{{.ID}}/bar-import-{{.ShortID}}.bat"><code>bar-import-{{.ShortID}}.bat</code></a>,
	save in root of the working tree and run.</p>
<p>This script is no-brain solution to import spec {{.ShortID}}. It automatically
download <code>barc.exe</code> if it is not found and download spec and BLOBs from bard</p>
<p>You can also download <a href="{{.Info.HTTPEndpoint}}/win/barc.exe"><code>barc.exe</code></a> and save
	it root of the working tree or somewhere in PATH.</p>
</body>
</html>
`

type SpecHandler struct {
	storage.Storage
	Info *proto.Info
	BarExe string
}

func (h *SpecHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	id := proto.ID(strings.TrimPrefix(r.URL.Path, "/v1/spec/"))

	logx.Debugf("serving spec %s", id)

	spec, err := h.Storage.IsSpecExists(id)
	if err != nil {
		logx.Errorf("bad spec id %s", id)
		w.WriteHeader(404)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	spec_template, err := h_template.New("spec").Parse(spec_tpl)
	if err != nil {
		logx.Error(err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err = spec_template.Execute(w, map[string]interface{}{
		"Info": h.Info,
		"BarExe": h.BarExe,
		"ID": id,
		"ShortID": id[:12],
	}); err != nil {
		logx.Error(err)
		w.WriteHeader(500)
		return
	}

	if err = json.NewEncoder(w).Encode(&spec); err != nil {
		logx.Error(err)
		w.WriteHeader(500)
		return
	}
}
