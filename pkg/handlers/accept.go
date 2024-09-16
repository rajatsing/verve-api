package handlers

import (
	"log"
	"net/http"
	"strconv"
	external "verve/pkg/external"
	uniqueids "verve/pkg/uniqueIds"
)

func AcceptHandler(w http.ResponseWriter, r *http.Request) {
	// Get the 'id' query parameter
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		// Return "failed" if 'id' is missing
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("failed"))
		return
	}

	// Parse the 'id' parameter to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// Return "failed" if 'id' is invalid
		respondFailedToWriter(w)
		return
	}

	// Add the id to the unique IDs set
	err = uniqueids.AddID(id)
	if err != nil {
		// Return "failed" if there's an error adding the ID
		respondFailedToWriter(w)
		return
	}

	// Get the 'endpoint' query parameter
	endpoint := r.URL.Query().Get("endpoint")

	if endpoint != "" {
		// Get the current unique ID count
		count := uniqueids.GetCurrentCount()

		// Fire an HTTP POST request to the provided endpoint with the count
		err = external.SendCountToEndpoint(endpoint, count)
		if err != nil {
			log.Printf("Error sending count to endpoint: %v", err)
			// Return "failed" if there's an error sending to the endpoint
			respondFailedToWriter(w)
			return
		}
	}

	// Return "ok" if no errors
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("ok"))
}

func respondFailedToWriter(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("failed"))
}
