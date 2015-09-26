package main

import (
	"encoding/base64"
	"io"
	"io/ioutil"
	"net/http"
)

// pasteHandler reads POST data from the file field and sets it in a cookie
// for the processor to take over.
func uploaderHandler(w http.ResponseWriter, req *http.Request) {
	// Use io.Reader type of req.FormFile to read the file and headers
	file, _, err := req.FormFile("textFile")
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	// Receive the bytes from the file and store them in a data variable
	data, err := ioutil.ReadAll(file)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "appText",
		Value: base64.StdEncoding.EncodeToString(data),
		Path:  "/",
	})

	w.Header().Set("Location", "/process")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
