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

// Chunker provides a skeleton for objects which can split data into smaller chunks of data
type Chunker interface {
	// Chunk takes in data and returns that data split into digestible
	// poritions to use for parallel processing
	Chunk() (map[int]*Chunk, error)
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
func (c SentenceChunker) Chunk() (map[int]*Chunk, error) {
	result := make(map[int]*Chunk)

	s := bufio.NewScanner(bytes.NewBufferString(c.Input))

	index := 0
	firstWord := ""
	for s.Scan() {
		for _, r := range s.Text() {
			if result[index] == nil {
				result[index] = NewChunk(index, "")
			}

			result[index].Data += string(r)

			if result[index].FirstWord == "" {
				if IsAlphaNumeric(r) {
					firstWord += string(r)
				} else {
					result[index].FirstWord = firstWord
					firstWord = ""
				}
			}
			if IsEndOfSentence(r) {
				index += 1
				continue
			}
		}
	}

	return result, nil
}
