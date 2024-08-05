package tpls

import (
	"encoding/json"
	"fmt"
	"ginapp/tpl"
	"html/template"
	"net/http"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
)

/*
{{ $data := . | structToMap }}
{{ range $key, $value := $data }}

	<tr>
	    <th scope="col">{{ $key }}</th><td> {{ $value }}</td>
	</tr>

{{ end }}
*/
func structToMap(item interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	v := reflect.ValueOf(item)

	for i := 0; i < v.NumField(); i++ {
		key := v.Type().Field(i).Name
		value := v.Field(i).Interface()
		out[key] = value
	}
	return out
}
func TplRouter(e *gin.Engine) {
	// 1. funcMap
	funcMap := template.FuncMap{
		"structToMap": structToMap,
		"formatAsDate": func(t time.Time) string {
			year, month, day := t.Date()
			return fmt.Sprintf("%d%02d/%02d", year, month, day)
		},
		"jsonEncode": func(obj interface{}) string {
			if d, err := json.Marshal(obj); err == nil {
				return string(d)
			} else {
				return err.Error()
			}
		},
	}
	// 2. config
	e.Delims("{{", "}}")

	/** drop
	e.SetFuncMap(funcMap)
	e.LoadHTMLGlob("./tpl/login/*")
	**/

	// 3. set template path with funcs
	tmplWithFuncs := template.Must(template.New("").Funcs(funcMap).ParseFS(tpl.GetLoginFS(), "login/*.tmpl"))
	e.SetHTMLTemplate(tmplWithFuncs)
	// e.LoadHTMLGlob("template/*/*.html")// 或者使用这种方式

	// 4. curl -D- m:4500/tpl-raw/login/redirect.tmpl
	e.StaticFS("tpl-raw", http.FS(tpl.GetLoginFS()))
	// 5. curl m:4500/tpl/redirect -d '{}'
	e.POST("tpl/redirect", TplRedirectPage)
}
