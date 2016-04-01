package main

import (
	"strings"
	"unicode"
)

const SentenceEnders = ".!?"

// IsAlpha checks if a current rune is a letter
func IsAlpha(r rune) bool {
	return unicode.IsLetter(r)
}

// IsSpace checks if a current rune is a space
func IsSpace(r rune) bool {
	return unicode.IsSpace(r)
}

// IsAlphaNumeric checks if a current rune is a letter or number
func IsAlphaNumeric(r rune) bool {
	return IsAlpha(r) || unicode.IsNumber(r)
}

// IsEndOfSentence checks if we have punctuation to end a sentence
func IsEndOfSentence(r rune) bool {
	return strings.ContainsRune(SentenceEnders, r)
}
