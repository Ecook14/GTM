package pagespeed

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from a .env file.
func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	return nil
}

// PageSpeedResponse represents the structure of the JSON response from the PageSpeed Insights API.
type PageSpeedResponse struct {
	Kind                 string `json:"kind"`
	ID                   string `json:"id"`
	LoadingExperience    struct {
		Metrics struct {
			FirstInputDelayP30 struct {
				Category string `json:"category"`
			} `json:"FIRST_INPUT_DELAY_P30"`
			// Add other metrics as needed
		} `json:"metrics"`
		// Add other fields as needed
	} `json:"loadingExperience"`
	LighthouseResult struct {
		Categories struct {
			Performance struct {
				Score float64 `json:"score"`
			} `json:"performance"`
			// Add other categories as needed
		} `json:"categories"`
		Audits struct {
			FirstContentfulPaint struct {
				DisplayValue string `json:"displayValue"`
				// Add other fields as needed
			} `json:"first-contentful-paint"`
			LargestContentfulPaint struct {
				DisplayValue string `json:"displayValue"`
				// Add other fields as needed
			} `json:"largest-contentful-paint"`
			TimeToInteractive struct {
				DisplayValue string `json:"displayValue"`
				// Add other fields as needed
			} `json:"time-to-interactive"`
			CumulativeLayoutShift struct {
				DisplayValue string `json:"displayValue"`
				// Add other fields as needed
			} `json:"cumulative-layout-shift"`
			// Add other audits as needed
		} `json:"audits"`
		// Add other fields as needed
		Version string `json:"version"`
		Report struct {
			URL      string `json:"url"`
			Waterfall []struct {
				StartTime float64 `json:"startTime"`
				EndTime    float64 `json:"endTime"`
				RequestID string `json:"requestId"`
				URL        string `json:"url"`
				ResponseReceivedTime float64 `json:"responseReceivedTime"`
				TimeToFirstByte      float64 `json:"timeToFirstByte"`
				// Add other fields as needed
			} `json:"waterfall"`
			Issues []struct {
				Code      string `json:"code"`
				Title     string `json:"title"`
				Description string `json:"description"`
				// Add other fields as needed
			} `json:"issues"`
		} `json:"report"`
	} `json:"lighthouseResult"`
	// Add other fields as needed
}

// GetPageSpeedScore retrieves the PageSpeed Score for a given URL.
func GetPageSpeedScore(url string) (string, error) {
	apiKey := os.Getenv("PSI_API_KEY")
	if apiKey == "" {
		return "", errors.New("API_KEY not set")
	}

	// Construct the API request URL
	apiURL := "https://www.googleapis.com/pagespeedonline/v5/runPagespeed?url=" + url + "&key=" + apiKey

	// Make the HTTP request
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to get PageSpeed Score")
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Parse the JSON response
	var result PageSpeedResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	// Extract the PageSpeed Score from the response
	score := result.LighthouseResult.Categories.Performance.Score

	// Extract other metrics and fields as needed
	// This is a placeholder; you'll need to adjust this based on the specific fields you're interested in
	output := fmt.Sprintf("PageSpeed Score: %.2f\n", score)

	// Add more detailed output based on the fields you've included in the PageSpeedResponse struct

	return output, nil
}
