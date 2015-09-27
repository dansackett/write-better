package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	proc "github.com/dansackett/go-text-processors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// Processor is an interface which handles processing of a Chunk.
type Processor interface {
	Process(*Chunk) *Chunk
}

// ActiveProcessors stores a list of processors that can be used to process on
// a given chunk. This provides an easy way to handle chaining of multiple
// processors.
type ActiveProcessors []Processor

// Process method satisfies the interface for a Processor allowing us to call
// this method on our list and run the chunk through each processor.
func (p ActiveProcessors) Process(c *Chunk) *Chunk {
	for _, processor := range p {
		c = processor.Process(c)
	}

	return c
}

// processorsHandler chunks and processes text
func processorsHandler(w http.ResponseWriter, r *http.Request) {
	data, err := getCookie(r, "appText")
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	var wg sync.WaitGroup

	// Chunk the text into sentences
	c := NewSentenceChunker(string(data))
	chunks, _ := c.Chunk()

	// Send each chunk into a gorountine to process
	for _, c := range chunks {
		wg.Add(1)
		go func(c *Chunk) {
			defer wg.Done()
			c = processors.Process(c)
		}(c)
	}

	// Wait for the processing to finish
	wg.Wait()

	// Format chunks map to be JSONified
	formattedChunks := make(map[string]*Chunk)
	for idx, chunk := range chunks {
		formattedChunks[strconv.Itoa(idx)] = chunk
	}

	jsonChunks, err := json.Marshal(formattedChunks)

	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "appResult",
		Value: base64.StdEncoding.EncodeToString(jsonChunks),
		Path:  "/results",
	})

	w.Header().Set("Location", "/results")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// doTextProcessor is a convenience function to make this more DRY. It runs
// the processors from go-text-processors giving the same treatment to each
// processor.
func doTextProcessor(p proc.TextProcessor, label string, c *Chunk, msg string) *Chunk {
	res := p.Run(c.Data)

	for _, match := range res.Matches {
		formattedMsg := fmt.Sprintf(msg, match.Match)
		c.Matches = append(c.Matches, NewMatch(match.Match, label, match.Indices, formattedMsg))
		c.Score += 1
	}

	return c
}

// PassiveVoiceProcessor is an empty struct for processing the passive voice
type PassiveVoiceProcessor struct{}

// UsePassiveVoiceProcessor is a convenience variable for referencing a PassiveVoiceProcessor
var UsePassiveVoiceProcessor PassiveVoiceProcessor

// Process handles the processing for passive voice matches
func (_ PassiveVoiceProcessor) Process(c *Chunk) *Chunk {
	msg := "\"%s\" is considered passive voice."
	return doTextProcessor(proc.PassiveVoiceProcessor(), "passive", c, msg)
}

// WeaselWordProcessor is an empty struct for processing weasel words
type WeaselWordProcessor struct{}

// UseWeaselWordProcessor is a convenience variable for referencing a WeaselWordProcessor
var UseWeaselWordProcessor WeaselWordProcessor

// Process handles the processing for weasel word matches
func (_ WeaselWordProcessor) Process(c *Chunk) *Chunk {
	msg := "\"%s\" is considered a weasel word."
	return doTextProcessor(proc.WeaselWordProcessor(), "weasel", c, msg)
}

// TooWordyProcessor is an empty struct for processing wordy phrases
type TooWordyProcessor struct{}

// UseTooWordyProcessor is a convenience variable for referencing a TooWordyProcessor
var UseTooWordyProcessor TooWordyProcessor

// Process handles the processing for wordy phrase matches
func (_ TooWordyProcessor) Process(c *Chunk) *Chunk {
	msg := "\"%s\" is considered a wordy phrase."
	return doTextProcessor(proc.TooWordyProcessor(), "wordy", c, msg)
}

// AdverbProcessor is an empty struct for processing adverbs
type AdverbProcessor struct{}

// UseAdverbProcessor is a convenience variable for referencing a AdverbProcessor
var UseAdverbProcessor AdverbProcessor

// Process handles the processing for adverb matches
func (_ AdverbProcessor) Process(c *Chunk) *Chunk {
	msg := "\"%s\" is an adverb."
	return doTextProcessor(proc.AdverbProcessor(), "adverb", c, msg)
}

// ClicheProcessor is an empty struct for processing cliches
type ClicheProcessor struct{}

// UseClicheProcessor is a convenience variable for referencing a ClicheProcessor
var UseClicheProcessor ClicheProcessor

// Process handles the processing for cliche matches
func (_ ClicheProcessor) Process(c *Chunk) *Chunk {
	msg := "\"%s\" is a cliche."
	return doTextProcessor(proc.ClicheProcessor(), "cliche", c, msg)
}

// LexicalIllusionProcessor is an empty struct for processing repeated words
type LexicalIllusionProcessor struct{}

// UseLexicalIllusionProcessor is a convenience variable for referencing a LexicalIllusionProcessor
var UseLexicalIllusionProcessor LexicalIllusionProcessor

// Process handles the processing for repeated word matches
func (_ LexicalIllusionProcessor) Process(c *Chunk) *Chunk {
	msg := "\"%s\" a repeated word."
	return doTextProcessor(proc.LexicalIllusionProcessor(), "lexical_illusion", c, msg)
}

// SentenceLengthProcessor is an empty struct for processing a sentence's length
type SentenceLengthProcessor struct{}

// UseSentenceLengthProcessor is a convenience variable for referencing a SentenceLengthProcessor
var UseSentenceLengthProcessor SentenceLengthProcessor

// Process handles the processing for long sentence matches
func (_ SentenceLengthProcessor) Process(c *Chunk) *Chunk {
	var indices []int

	if len(c.Data) > 160 {
		msg := fmt.Sprintf("This is a VERY long sentence.")
		c.Matches = append(c.Matches, NewMatch("", "length", indices, msg))
		c.Score += 1
	} else if len(c.Data) > 130 {
		msg := fmt.Sprintf("This is a long sentence.")
		c.Matches = append(c.Matches, NewMatch("", "length", indices, msg))
		c.Score += 1
	}

	return c
}

// StartsWithProcessor is an empty struct for processing a sentence's first phrase
type StartsWithProcessor struct{}

// UseStartsWithProcessor is a convenience variable for referencing a StartsWithProcessor
var UseStartsWithProcessor StartsWithProcessor

// Process handles the processing for first phrase matches
func (_ StartsWithProcessor) Process(c *Chunk) *Chunk {
	// @TODO Check if it starts with "there is" or "there are"
	var indices []int

	if strings.ToLower(c.Data) == "so" {
		msg := fmt.Sprintf("This sentence starts with so. Consider changing it.")
		c.Matches = append(c.Matches, NewMatch("", "length", indices, msg))
		c.Score += 1
	}

	return c
}
