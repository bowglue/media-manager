package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"app/cmd/videostream/config"
	"app/cmd/videostream/handlers"
)

func main() {
	// Create output directory
	os.MkdirAll(config.OutputDir, 0755)

	// Create HTTP handlers for streaming
	http.HandleFunc("/stream", handlers.StreamHandler)
	http.HandleFunc("/segments/", handlers.SegmentsHandler)

	// Start the server
	fmt.Printf("Streaming server started on port %d\n", config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil))
}