// pkg/external/external.go

// Package external contains functions to interact with external services.
package external

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// CountData represents the data to be sent in the POST request.
type CountData struct {
	Count int `json:"count"`
}

// SendCountToEndpoint sends the count to the provided endpoint via POST request.
// Logs the HTTP status code of the response.
func SendCountToEndpoint(endpoint string, count int) error {
	// Create the JSON payload
	data := CountData{Count: count}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Create a new POST request
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	// Set up the HTTP client with a timeout
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Log the HTTP status code of the response
	log.Printf("POST request to %s returned status code %d", endpoint, resp.StatusCode)

	return nil
}
