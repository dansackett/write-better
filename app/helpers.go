package main

import (
	"encoding/base64"
	"net/http"
	"strings"
)

const (
	SentenceEnders = ".!?"
	AlphaNumeric   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
)

// IsAlphaNumeric checks if a current rune is a letter or number
func IsAlphaNumeric(r rune) bool {
	return strings.ContainsRune(AlphaNumeric, r)
}

// IsEndOfSentence checks if we have punctuation to end a sentence
func IsEndOfSentence(r rune) bool {
	return strings.ContainsRune(SentenceEnders, r)
}

// getCookie is a convenience method for fetching and decoding a cookie
func getCookie(req *http.Request, name string) ([]byte, error) {
	res, err := req.Cookie(name)

	if err != nil {
		return nil, err
	}

	data, err := base64.StdEncoding.DecodeString(res.Value)

	if err != nil {
		return nil, err
	}

	return data, nil
}
