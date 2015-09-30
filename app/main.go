package main

import (
	"fmt"
	"log"
	"net/http"
)

// processors refer to the list of active processors we are running in the
// application at any given time
var processors Processor = ActiveProcessors{
	UsePassiveVoiceProcessor,
	UseWeaselWordProcessor,
	UseTooWordyProcessor,
	UseAdverbProcessor,
	UseClicheProcessor,
	UseLexicalIllusionProcessor,
	UseSentenceLengthProcessor,
	UseStartsWithProcessor,
	// This needs to be the last one since it manipuates the string itself
	UseHTMLProcessor,
}

// appText refers to the text that has been submitted
var appText string

// appResult refers to the processing results
var appResult Chunks

func main() {
	http.Handle("/", &templateHandler{filename: "index.html"})
	http.HandleFunc("/upload", uploaderHandler)
	http.HandleFunc("/paste", pasteHandler)
	http.HandleFunc("/process", processorsHandler)
	http.HandleFunc("/results", resultHandler)

	fmt.Println("App server running on :8080")

	// Start the web server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
