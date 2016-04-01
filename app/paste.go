package main

import (
	"bytes"
	"net/http"
)

// pasteHandler reads POST data from the textarea field and sets it in a
// cookie for the processor to take over.
func pasteHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	data := req.Form.Get("textFile")

	// This is hackish as it replaces a character I am having difficulty with.
	cleanedData := bytes.Replace([]byte(data), []byte("’"), []byte("'"), -1)
	cleanedData = bytes.Replace([]byte(cleanedData), []byte("‘"), []byte("'"), -1)
	cleanedData = bytes.Replace([]byte(cleanedData), []byte("“"), []byte("\""), -1)
	cleanedData = bytes.Replace([]byte(cleanedData), []byte("”"), []byte("\""), -1)

	appText = string(cleanedData)

	w.Header().Set("Location", "/process")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
