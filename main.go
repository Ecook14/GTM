package main

import (
	"log"
	"net/http"

	"github.com/Ecook14/GoTM/api"
)

func main() {
	http.HandleFunc("/analyze", api.HandleAnalysisRequest)

	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
