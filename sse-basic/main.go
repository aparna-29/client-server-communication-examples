package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func sseHandler(w http.ResponseWriter, r *http.Request) {
	// Set required headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Optional: for CORS if needed

	// Flush support
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Example loop to send events every second
	for i := 0; i < 10; i++ {
		fmt.Fprintf(w, "data: Message number %d at %s\n\n", i, time.Now().Format(time.RFC3339))
		flusher.Flush()

		time.Sleep(1 * time.Second)
	}
}

func main() {
	http.HandleFunc("/events", sseHandler)

	fmt.Println("Starting server at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
