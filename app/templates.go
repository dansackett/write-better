package main

import (
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

// templateHandler allows us to load an HTML file and serve it. We parse the
// template once so we don't waste resources loading it over and over.
type templateHandler struct {
	Once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.Once.Do(func() {
		path, _ := filepath.Abs(filepath.Join("templates", t.filename))
		t.templ = template.Must(template.ParseFiles(path))
	})

	var data map[string]interface{}

	t.templ.Execute(w, data)
}
