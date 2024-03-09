package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Ecook14/GoTM/external/pagespeed"
)

func main() {
	// Load environment variables
	if err := pagespeed.LoadEnv(); err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}

	// Define the HTTP handler function
	http.HandleFunc("/pagespeed", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if url == "" {
			http.Error(w, "Missing 'url' query parameter", http.StatusBadRequest)
			return
		}

		// Fetch the PageSpeed Score and other metrics
		score, err := pagespeed.GetPageSpeedScore(url)
		if err != nil {
			http.Error(w, "Failed to fetch PageSpeed Score", http.StatusInternalServerError)
			return
		}

		// Send the response
		fmt.Fprintf(w, "PageSpeed Score: %s\n", score)
	})

	// Start the HTTP server
	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
