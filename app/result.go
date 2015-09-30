package main

import (
	"net/http"
	"sort"
	"strings"
)

// resultHandler reads the result cookie, parses it, and gets it ready to be
// used in a template to show users how their text finishes.
func resultHandler(w http.ResponseWriter, req *http.Request) {
	var fullText []string
	var score int

	chunks := appResult

	sort.Sort(ByChunk(chunks))

	for _, chunk := range chunks {
		fullText = append(fullText, chunk.Data)
		score += chunk.Score
	}

	// Build data for the template
	returnData := map[string]interface{}{
		"fullText": strings.Join(fullText, ""),
		"score":    interface{}(score),
	}

	// Render the template
	t := &templateHandler{filename: "results.html"}
	RenderTemplate(t, w, returnData)
}
