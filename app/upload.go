package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

// uploadHandler reads POST data from the file field and sets it in the app
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

	// This is hackish as it replaces a character I am having difficulty with.
	data = bytes.Replace([]byte(data), []byte("’"), []byte("'"), -1)
	data = bytes.Replace([]byte(data), []byte("‘"), []byte("'"), -1)
	data = bytes.Replace([]byte(data), []byte("“"), []byte("\""), -1)
	data = bytes.Replace([]byte(data), []byte("”"), []byte("\""), -1)

	appText = string(data)

	w.Header().Set("Location", "/process")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
