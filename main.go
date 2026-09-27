package main

import (
	"fmt"
	"net/http"
)

// startServer creates the HTTP server and registers all routes.
func startServer() {

	// Create a new router for handling HTTP requests.
	mux := http.NewServeMux()

	// Root endpoint: GET /
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "URL Shortener")
	})

	// Health endpoint: used to check whether the server is running.
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	// Shorten endpoint: will be used to create a short URL.
	mux.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Shorten endpoint")
	})

	// Start the server on port 8080.
	fmt.Println("Server running on http://localhost:8080")

	err := http.ListenAndServe(":8080", mux)

	// Print an error if the server stops unexpectedly.
	if err != nil {
		fmt.Println("Server error:", err)
	}
}

func main() {

	// Start the HTTP server.
	startServer()
}
