package main

import (
	"fmt"
	"log"
	"net/http"
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

func main() {
	http.Handle("/", &templateHandler{filename: "index.html"})
	http.HandleFunc("/upload", uploaderHandler)
	http.HandleFunc("/paste", pasteHandler)
	http.HandleFunc("/process", processorsHandler)

	fmt.Println("App server running on :8080")

	// Start the web server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
