package main

import (
	"fmt"
	"sync"
	// "log"
	// "net/http"
)

var processors Processor = ActiveProcessors{
	UsePassiveVoiceProcessor,
	UseWeaselWordProcessor,
	UseTooWordyProcessor,
	UseAdverbProcessor,
	UseClicheProcessor,
	UseLexicalIllusionProcessor,
	UseSentenceLengthProcessor,
	UseStartsWithProcessor,
}

var wg sync.WaitGroup

func main() {
	// Exceprt from http://microfictionmondaymagazine.com/ by Pavelle Wesser
	s := `I was on fire after winning the science competition, which may be
	why, as I was accepting the trophy, it disintegrated in my hands while my
	synapses short-circuited. Through the haze of my mind, I tried to tell Dad
	the pics he was snapping of me would be his last. \"Dad!\" The word burned
	to cinders before emerging from my charred lips. I extended my arms, which
	exploded off my shoulders, prompting piercing screams from the audience.
	Finally, I combusted, and the immense pressure that had been building up
	within me from the beginning of the competition was released.`

	// Chunk the text into sentences
	c := NewSentenceChunker(s)
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

	// Print the results of the processing
	for _, c := range chunks {
		fmt.Printf("%d. (Score: %d)\n", c.Index, c.Score)
		fmt.Println("---")

		if len(c.Messages) > 0 {
			for _, msg := range c.Messages {
				fmt.Println("-", msg)
			}
		} else {
			fmt.Println("- No errors to note")
		}
	}
}
