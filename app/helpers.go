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
