package main

import (
	"net/http"

	unidecode "github.com/rainycape/unidecode"
)

// pasteHandler reads POST data from the textarea field and sets it in a
// cookie for the processor to take over.
func pasteHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	data := req.Form.Get("textFile")
	appText = unidecode.Unidecode(string(data))

	w.Header().Set("Location", "/process")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
