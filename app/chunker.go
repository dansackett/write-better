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
	// Messages store the helpful messages returned from processors
	Matches []*Match
	// Score refers to the overall score of this Chunk
	Score int
}

// NewChunk is a convenience function to build a new Chunk instance
func NewChunk(idx int, data string) *Chunk {
	var matches []*Match

	return &Chunk{
		Index:     idx,
		Data:      data,
		FirstWord: "",
		Matches:   matches,
		Score:     0,
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

	s := bufio.NewScanner(bytes.NewBufferString(c.Input))

	index := 0
	firstWord := ""
	for s.Scan() {
		for _, r := range s.Text() {
			if tmp[index] == nil {
				tmp[index] = NewChunk(index, "")
			}

			tmp[index].Data += string(r)

			if tmp[index].FirstWord == "" {
				if IsAlphaNumeric(r) {
					firstWord += string(r)
				} else {
					tmp[index].FirstWord = firstWord
					firstWord = ""
				}
			}
			if IsEndOfSentence(r) {
				index += 1
				continue
			}
		}
	}

	// Assign chunks
	for _, chunk := range tmp {
		result = append(result, chunk)
	}

	return result, nil
}
