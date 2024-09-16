package main

import (
	"log"
	"net/http"
	"os"
	"verve/pkg/handlers"
)

func main() {

	// Get the port from the environment variable, default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// route setup
	http.HandleFunc("/api/verve/accept", handlers.AcceptHandler)

	// start the server
	log.Printf("server starting on port %s...", port)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
