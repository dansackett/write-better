package main

import (
	"bytes"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

// Average words per minute
const AvgReadingSpeed = 275

// resultHandler reads the result cookie, parses it, and gets it ready to be
// used in a template to show users how their text finishes.
func resultHandler(w http.ResponseWriter, req *http.Request) {
	var score int
	var fullText []string
	var curStr bytes.Buffer

	chunks := appResult
	matches := map[string]int{
		"passive":    0,
		"weasel":     0,
		"wordy":      0,
		"adverb":     0,
		"cliche":     0,
		"illusion":   0,
		"length":     0,
		"startswith": 0,
	}

	sort.Sort(ByChunk(chunks))

	for _, chunk := range chunks {
		// Build paragraphs
		if chunk.IsNewParagraph {
			curStr.WriteString("</p>")
			fullText = append(fullText, curStr.String())

			curStr.Reset()
			curStr.WriteString("<p>")
			curStr.WriteString(chunk.Data)
		} else {
			curStr.WriteString(chunk.Data)
		}

		// Build match sums
		if len(chunk.Matches) > 0 {
			for _, match := range chunk.Matches {
				matches[match.Label] += 1
			}
		}

		score += chunk.Score
	}

	// Build data for the template
	returnData := map[string]interface{}{
		"score":    score,
		"matches":  matches,
		"summary":  appSummary,
		"readTime": GetReadTime(),
		"fullText": fullText,
	}

	// Render the template
	t := &templateHandler{filename: "results.html"}
	RenderTemplate(t, w, returnData)
}

func GetReadTime() string {
	var buffer bytes.Buffer

	readTime := float64(appSummary["words"]) / float64(AvgReadingSpeed)
	vals := strings.Split(strconv.FormatFloat(readTime, 'f', 2, 64), ".")
	beforeDecimal, _ := strconv.Atoi(vals[0])
	afterDecimal, _ := strconv.ParseFloat(vals[1], 64)

	hours := beforeDecimal / 60
	minutes := beforeDecimal % 60
	seconds := afterDecimal * .6

	if hours > 0 {
		buffer.WriteString(strconv.Itoa(hours))
		if hours == 1 {
			buffer.WriteString(" hour ")
		} else {
			buffer.WriteString(" hours ")
		}
	}

	if minutes > 0 {
		buffer.WriteString(strconv.Itoa(minutes))
		if minutes == 1 {
			buffer.WriteString(" minute ")
		} else {
			buffer.WriteString(" minutes ")
		}
	}

	if seconds > 0 {
		vals := strings.Split(strconv.FormatFloat(seconds, 'f', 2, 64), ".")
		buffer.WriteString(vals[0])
		if seconds == 1 {
			buffer.WriteString(" second ")
		} else {
			buffer.WriteString(" seconds")
		}
	}

	return buffer.String()
}
