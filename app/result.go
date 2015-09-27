package main

import (
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

	// Get Result data
	data, err := getCookie(req, "appResult")
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	// Convert JSON result to Golang object
	err = json.Unmarshal(data, &chunks)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	// Get fullText used for parsing
	fullText, err := getCookie(req, "appText")
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	// Build data for the template
	returnData := map[string]interface{}{
		"fullText": string(fullText),
		"json":     string(data),
	}

	// Calculate the final score of the writing
	score := 0
	for _, chunk := range chunks {
		score += chunk.Score
	}
	returnData["score"] = interface{}(score)

	// Render the template
	t := &templateHandler{filename: "results.html"}
	t.Once.Do(func() {
		path, _ := filepath.Abs(filepath.Join("templates", t.filename))
		t.templ = template.Must(template.ParseFiles(path))
	})
	t.templ.Execute(w, returnData)
}
