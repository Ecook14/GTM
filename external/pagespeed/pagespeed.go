package pagespeed

import (
	"encoding/json"
	"errors"
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
	LighthouseResult struct {
		Categories struct {
			Performance struct {
				Score float64 `json:"score"`
			} `json:"performance"`
		} `json:"categories"`
		Audits struct {
			FirstContentfulPaint struct {
				DisplayValue string `json:"displayValue"`
			} `json:"first-contentful-paint"`
			LargestContentfulPaint struct {
				DisplayValue string `json:"displayValue"`
			} `json:"largest-contentful-paint"`
			TimeToInteractive struct {
				DisplayValue string `json:"displayValue"`
			} `json:"time-to-interactive"`
			CumulativeLayoutShift struct {
				DisplayValue string `json:"displayValue"`
			} `json:"cumulative-layout-shift"`
		} `json:"audits"`
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
			} `json:"waterfall"`
			Issues []struct {
				Code      string `json:"code"`
				Title     string `json:"title"`
				Description string `json:"description"`
			} `json:"issues"`
		} `json:"report"`
	} `json:"lighthouseResult"`
}

// GetPageSpeedScore retrieves the PageSpeed Score and detailed metrics for a given URL.
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

	// Construct a detailed response
	detailedResponse := map[string]interface{}{
		"PageSpeed Score": result.LighthouseResult.Categories.Performance.Score,
		"Performance Metrics": map[string]interface{}{
			"First Contentful Paint": result.LighthouseResult.Audits.FirstContentfulPaint.DisplayValue,
			"Largest Contentful Paint": result.LighthouseResult.Audits.LargestContentfulPaint.DisplayValue,
			"Time to Interactive": result.LighthouseResult.Audits.TimeToInteractive.DisplayValue,
			"Cumulative Layout Shift": result.LighthouseResult.Audits.CumulativeLayoutShift.DisplayValue,
		},
		"Waterfall": result.LighthouseResult.Report.Waterfall,
		"Issues": result.LighthouseResult.Report.Issues,
	}

	// Convert the detailed response to a JSON string
	jsonResponse, err := json.MarshalIndent(detailedResponse, "", " ")
	if err != nil {
		return "", err
	}

	return string(jsonResponse), nil
}
