package main

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"text/template"
)

// resultHandler reads the result cookie, parses it, and gets it ready to be
// used in a template to show users how their text finishes.
func resultHandler(w http.ResponseWriter, req *http.Request) {
	chunks := make(map[string]*Chunk)
	res, err := req.Cookie("appResult")

	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	data, err := base64.StdEncoding.DecodeString(res.Value)

	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	err = json.Unmarshal(data, &chunks)

	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	t := &templateHandler{filename: "results.html"}

	t.Once.Do(func() {
		path, _ := filepath.Abs(filepath.Join("templates", t.filename))
		t.templ = template.Must(template.ParseFiles(path))
	})

	t.templ.Execute(w, chunks)
}
