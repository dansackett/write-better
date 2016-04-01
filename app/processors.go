package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	proc "github.com/dansackett/go-text-processors"
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
	var wg sync.WaitGroup

	// Chunk the text into sentences
	c := NewSentenceChunker(appText)
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

	appResult = chunks

	w.Header().Set("Location", "/results")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// doTextProcessor is a convenience function to make this more DRY. It runs
// the processors from go-text-processors giving the same treatment to each
// processor.
func doTextProcessor(p proc.TextProcessor, label string, c *Chunk, msg string) *Chunk {
	res := p.Run(c.Data)

	for _, match := range res.Matches {
		formattedMsg := fmt.Sprintf(msg)
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
	msg := "This is considered passive voice."
	return doTextProcessor(proc.PassiveVoiceProcessor(), "passive", c, msg)
}

// WeaselWordProcessor is an empty struct for processing weasel words
type WeaselWordProcessor struct{}

// UseWeaselWordProcessor is a convenience variable for referencing a WeaselWordProcessor
var UseWeaselWordProcessor WeaselWordProcessor

// Process handles the processing for weasel word matches
func (_ WeaselWordProcessor) Process(c *Chunk) *Chunk {
	msg := "This is considered a weasel word."
	return doTextProcessor(proc.WeaselWordProcessor(), "weasel", c, msg)
}

// TooWordyProcessor is an empty struct for processing wordy phrases
type TooWordyProcessor struct{}

// UseTooWordyProcessor is a convenience variable for referencing a TooWordyProcessor
var UseTooWordyProcessor TooWordyProcessor

// Process handles the processing for wordy phrase matches
func (_ TooWordyProcessor) Process(c *Chunk) *Chunk {
	msg := "This is considered a wordy phrase."
	return doTextProcessor(proc.TooWordyProcessor(), "wordy", c, msg)
}

// AdverbProcessor is an empty struct for processing adverbs
type AdverbProcessor struct{}

// UseAdverbProcessor is a convenience variable for referencing a AdverbProcessor
var UseAdverbProcessor AdverbProcessor

// Process handles the processing for adverb matches
func (_ AdverbProcessor) Process(c *Chunk) *Chunk {
	msg := "This is an adverb."
	return doTextProcessor(proc.AdverbProcessor(), "adverb", c, msg)
}

// ClicheProcessor is an empty struct for processing cliches
type ClicheProcessor struct{}

// UseClicheProcessor is a convenience variable for referencing a ClicheProcessor
var UseClicheProcessor ClicheProcessor

// Process handles the processing for cliche matches
func (_ ClicheProcessor) Process(c *Chunk) *Chunk {
	msg := "This is a cliche."
	return doTextProcessor(proc.ClicheProcessor(), "cliche", c, msg)
}

// LexicalIllusionProcessor is an empty struct for processing repeated words
type LexicalIllusionProcessor struct{}

// UseLexicalIllusionProcessor is a convenience variable for referencing a LexicalIllusionProcessor
var UseLexicalIllusionProcessor LexicalIllusionProcessor

// Process handles the processing for repeated word matches
func (_ LexicalIllusionProcessor) Process(c *Chunk) *Chunk {
	msg := "This a repeated word."
	return doTextProcessor(proc.LexicalIllusionProcessor(), "illusion", c, msg)
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
	if strings.ToLower(c.FirstWord) == "so" {
		msg := fmt.Sprintf("This sentence starts with 'so'. Consider changing it.")
		c.Matches = append(c.Matches, NewMatch("", "startswith", getStartsWithIndices("s", 2, c), msg))
		c.Score += 1
	} else if strings.ToLower(c.FirstWord) == "there" {
		if strings.HasPrefix(strings.ToLower(c.Data), "there is") {
			msg := fmt.Sprintf("This sentence starts with 'there is'. Consider changing it.")
			c.Matches = append(c.Matches, NewMatch("", "startswith", getStartsWithIndices("t", 8, c), msg))
			c.Score += 1
		} else if strings.HasPrefix(strings.ToLower(c.Data), "there are") {
			msg := fmt.Sprintf("This sentence starts with 'there are'. Consider changing it.")
			c.Matches = append(c.Matches, NewMatch("", "startswith", getStartsWithIndices("t", 9, c), msg))
			c.Score += 1
		}
	}

	return c
}

// getStartsWithIndices helps find the correct indices for a starting phrase
// in cases the string begins with quotes or other characters.
func getStartsWithIndices(str string, strSize int, c *Chunk) []int {
	if strings.ToLower(string(c.Data[0])) == str {
		return []int{0, strSize}
	}
	firstOcc := 1
	for i, s := range c.Data {
		if strings.ToLower(string(s)) == str {
			firstOcc = i
			break
		}
	}
	return []int{firstOcc, firstOcc + strSize}
}

// HTMLProcessor applies HTML tags to the sentence for the frontend
type HTMLProcessor struct{}

// UseHTMLProcessor is a convenience variable for referencing a HTMLProcessor
var UseHTMLProcessor HTMLProcessor

// Process applies the HTML tags to the string
func (_ HTMLProcessor) Process(c *Chunk) *Chunk {
	nodes := ToCharNodes(c.Data)
	nodesLen := len(nodes)

	if len(c.Matches) > 0 {
		for _, match := range c.Matches {
			if len(match.Indices) == 0 {
				nodes[0].AddBefore(OpenTag(match.Label, match.Message))
				nodes[nodesLen-1].AddAfter(CloseTag())
			} else {
				nodes[match.Indices[0]].AddBefore(OpenTag(match.Label, match.Message))
				nodes[match.Indices[1]-1].AddAfter(CloseTag())
			}
		}
	}

	c.Data = nodes.ToString()

	return c
}
