package main

import (
	"encoding/base64"
	"net/http"
)

// pasteHandler reads POST data from the textarea field and sets it in a
// cookie for the processor to take over.
func pasteHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	http.SetCookie(w, &http.Cookie{
		Name:  "appText",
		Value: base64.StdEncoding.EncodeToString([]byte(req.Form.Get("textFile"))),
		Path:  "/",
	})

	w.Header().Set("Location", "/process")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
