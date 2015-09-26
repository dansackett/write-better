package main

import (
	"net/http"
)

// resultHandler reads the result cookie, parses it, and gets it ready to be
// used in a template to show users how their text finishes.
func resultHandler(w http.ResponseWriter, req *http.Request) {
	// @TODO read result cookie
	// @TODO Place into template
}
