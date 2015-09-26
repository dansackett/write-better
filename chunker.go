package main

import (
	"bufio"
	"bytes"
)

// Chunk is a piece of data created from the Chunker
type Chunk struct {
	// Index refers to the order of the Chunk
	Index int
	// Data is the text being stored
	Data string
	// FirstWord saves the first word of the text for analysis
	FirstWord string
	// Messages store the helpful messages returned from processors
	Messages []string
	// Errors store errors that have occurred
	Errors []string
	// Score refers to the overall score of this Chunk
	Score int
}

// NewChunk is a convenience function to build a new Chunk instance
func NewChunk(idx int, data string) *Chunk {
	var messages []string
	var errors []string

	return &Chunk{
		Index:     idx,
		Data:      data,
		FirstWord: "",
		Messages:  messages,
		Errors:    errors,
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
