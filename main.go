package main

import (
	"encoding/json"
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

	// Shorten endpoint: accepts a URL and returns a JSON response.
	mux.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {

		// Only allow POST requests.
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Request body structure.
		var request struct {
			URL string `json:"url"`
		}

		// Decode JSON request body.
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Check if URL was provided.
		if request.URL == "" {
			http.Error(w, "URL is required", http.StatusBadRequest)
			return
		}

		// Temporary response.
		response := struct {
			Message string `json:"message"`
			URL     string `json:"url"`
		}{
			Message: "URL received successfully",
			URL:     request.URL,
		}

		// Tell client that response is JSON.
		w.Header().Set("Content-Type", "application/json")

		// Send JSON response.
		json.NewEncoder(w).Encode(response)
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
