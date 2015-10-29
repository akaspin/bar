package front

import (
	h_template "html/template"
	"io"
	"strings"
	"text/template"
)

// weird. needs refactoring
type HttpTpl struct {
	tpls map[string]func(data map[string]interface{}, w io.Writer) (err error)
}

func NewHttpTpls() (res *HttpTpl, err error) {
	res = &HttpTpl{map[string]func(data map[string]interface{}, w io.Writer) (err error){}}

	h_funcmap := h_template.FuncMap{"join": strings.Join}
	funcmap := template.FuncMap{"join": strings.Join}

	// /
	frontTpl, err := h_template.New("front").Funcs(h_funcmap).Parse(front_tpl)
	if err != nil {
		return
	}
	res.tpls["front"] = func(data map[string]interface{}, w io.Writer) (err error) {
		err = frontTpl.Execute(w, data)
		return
	}

	// /spec
	specTpl, err := h_template.New("spec").Funcs(h_funcmap).Parse(spec_tpl)
	if err != nil {
		return
	}
	res.tpls["spec"] = func(data map[string]interface{}, w io.Writer) (err error) {
		err = specTpl.Execute(w, data)
		return
	}

	// /bats
	exportBatT, err := template.New("bar-export.bat").Funcs(funcmap).Parse(export_bat_tpl)
	if err != nil {
		return
	}
	res.tpls["bar-export.bat"] = func(data map[string]interface{}, w io.Writer) (err error) {
		err = exportBatT.Execute(w, data)
		return
	}
	importBatT, err := template.New("bar-import.bat").Funcs(funcmap).Parse(import_bat_tpl)
	if err != nil {
		return
	}
	res.tpls["bar-import.bat"] = func(data map[string]interface{}, w io.Writer) (err error) {
		err = importBatT.Execute(w, data)
		return
	}
	return
}
