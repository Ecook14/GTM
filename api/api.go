package api

import (
	"encoding/json"
	"net/http"

	"github.com/Ecook14/GTM/loader"
	"github.com/Ecook14/GTM/analyzer"
	"github.com/Ecook14/GTM/report"
)

type AnalysisResponse struct {
	URL      string `json:"url"`
	FCP      string `json:"fcp"`
	LCP      string `json:"lcp"`
	Analysis string `json:"analysis"`
}

func HandleAnalysisRequest(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "Missing 'url' query parameter", http.StatusBadRequest)
		return
	}

	// Load the page and capture metrics
	fcp, lcp, err := loader.LoadPageAndCaptureMetrics(url)
	if err != nil {
		http.Error(w, "Failed to load page and capture metrics", http.StatusInternalServerError)
		return
	}

	// Analyze the metrics
	analysis := analyser.AnalyzeMetrics(fcp, lcp)

	// Generate the report
	report := report.GenerateReport(analysis)

	// Prepare the response
	response := AnalysisResponse{
		URL:      url,
		FCP:      fcp.String(),
		LCP:      lcp.String(),
		Analysis: report,
	}

	// Set the response header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Write the response
	json.NewEncoder(w).Encode(response)
}
