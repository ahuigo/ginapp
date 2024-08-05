package tpls

import (
	"ginapp/tpl"
	"html/template"
	"log"
	"net/http"
	"testing"
)

var (
	pages = map[string]string{
		"/embed/index": "resource/index.tmpl",
	}
)

// curl m:4501/embed/index
func TmpEmbedServer(t *testing.T) {
	resourceFS := tpl.GetResourceFS()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		page, ok := pages[r.URL.Path]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		tpl, err := template.ParseFS(resourceFS, page)
		if err != nil {
			log.Printf("resource %s not found for uri %s", page, r.RequestURI)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		data := map[string]interface{}{
			"userAgent": r.UserAgent(),
		}
		if err := tpl.Execute(w, data); err != nil {
			return
		}
	})
	http.FileServer(http.FS(resourceFS))
	log.Println("server started...")
	err := http.ListenAndServe(":4501", nil)
	if err != nil {
		panic(err)
	}
}
