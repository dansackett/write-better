package main

import (
	"net/http"
)

// pasteHandler reads POST data from the textarea field and sets it in a
// cookie for the processor to take over.
func pasteHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	appText = string(req.Form.Get("textFile"))

	w.Header().Set("Location", "/process")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
