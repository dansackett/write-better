package main

import (
	"bufio"
	"bytes"
)

// Match represents an actual matched result after processing
type Match struct {
	// Match is the actual word / phrase that matches
	Match string
	// Label is the type of processor
	Label string
	// Indices are the start and end points of the match
	Indices []int
	// Message is the message from the processor
	Message string
}

// NewMatch is a convenience function to build a new Match instance
func NewMatch(match string, label string, indices []int, msg string) *Match {
	return &Match{
		Match:   match,
		Label:   label,
		Indices: indices,
		Message: msg,
	}
}

// Chunk is a piece of data created from the Chunker
type Chunk struct {
	// Index refers to the order of the Chunk
	Index int
	// Data is the text being stored
	Data string
	// FirstWord saves the first word of the text for analysis
	FirstWord string
	// IsNewParagraph marks when new paragraph delimiters are needed
	IsNewParagraph bool
	// Messages store the helpful messages returned from processors
	Matches []*Match
	// Score refers to the overall score of this Chunk
	Score int
}

// NewChunk is a convenience function to build a new Chunk instance
func NewChunk(idx int, data string) *Chunk {
	var matches []*Match

	return &Chunk{
		Index:          idx,
		Data:           data,
		FirstWord:      "",
		IsNewParagraph: false,
		Matches:        matches,
		Score:          0,
	}
}

// Chunks is a simple way to reference data that has been split by a chunker
type Chunks []*Chunk

// ByChunk is a sorting mechanism for sorting a slice of chunks
type ByChunk []*Chunk

func (c ByChunk) Len() int           { return len(c) }
func (c ByChunk) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c ByChunk) Less(i, j int) bool { return c[i].Index < c[j].Index }

// Chunker provides a skeleton for objects which can split data into smaller chunks of data
type Chunker interface {
	// Chunk takes in data and returns that data split into digestible
	// poritions to use for parallel processing
	Chunk() (Chunks, error)
}

// SentenceChunker is a Chunker instance which splits a string by sentences
type SentenceChunker struct {
	// Input is the text to be chunked
	Input string
}

// NewSentenceChunker is a convenience function to give us a SentenceChunker object
func NewSentenceChunker(input string) *SentenceChunker {
	return &SentenceChunker{Input: input}
}

// SentenceChunker takes the passed in input and splits it by sentences
func (c SentenceChunker) Chunk() (Chunks, error) {
	var result Chunks
	tmp := make(map[int]*Chunk)

	appSummary = map[string]int{
		"paragraphs": 0,
		"sentences":  0,
		"words":      0,
		"characters": 0,
		"letters":    0,
	}

	s := bufio.NewScanner(bytes.NewBufferString(c.Input))

	index := 0
	firstWord := ""
	for s.Scan() {
		text := s.Text()
		textLen := len(text)
		newPara := true

		var prevRune rune
		for i, r := range text {
			// Check if we have a new sentence.
			if tmp[index] == nil {
				tmp[index] = NewChunk(index, "")
				appSummary["sentences"] += 1

				// In the case of a new paragraph we add one for the first
				// word and increase the paragraph count
				if newPara {
					tmp[index].IsNewParagraph = true
					appSummary["paragraphs"] += 1
					appSummary["words"] += 1
				}

				newPara = false
			}

			// Add the current character to the chunk
			tmp[index].Data += string(r)

			// Increase the number of characters and possibly letters
			appSummary["characters"] += 1
			if IsAlpha(r) {
				appSummary["letters"] += 1
			}

			// Increase the word count
			if IsSpace(r) && !IsSpace(prevRune) {
				appSummary["words"] += 1
			}

			// Build the first word for easy reference later
			if tmp[index].FirstWord == "" {
				if IsAlphaNumeric(r) {
					firstWord += string(r)
				} else {
					tmp[index].FirstWord = firstWord
					firstWord = ""
				}
			}

			// We move on if we're at the end of the sentence or in the case
			// that a new paragraph does not have sentence terminators then we
			// must increase as well to keep paragraphs correct.
			if IsEndOfSentence(r) || textLen-1 == i {
				index++
			} else if textLen-1 == i {
				appSummary["words"] += 1
				index++
			}

			prevRune = r
		}
	}

	// Assign chunks
	for _, chunk := range tmp {
		result = append(result, chunk)
	}

	return result, nil
}
