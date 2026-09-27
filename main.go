package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
)

// generateShortCode creates a random 6-character short code.
func generateShortCode() string {

	const characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	code := make([]byte, 6)

	for i := range code {
		code[i] = characters[rand.Intn(len(characters))]
	}

	return string(code)
}

// startServer creates the HTTP server and registers all routes.
func startServer() {
	// Stores short code -> original URL.
	urls := make(map[string]string)
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
		// Generate a unique short code.
		shortCode := generateShortCode()

		// Store the URL using the short code.
		urls[shortCode] = request.URL

		response := struct {
			Message   string `json:"message"`
			ShortCode string `json:"short_code"`
			URL       string `json:"url"`
		}{
			Message:   "URL shortened successfully",
			ShortCode: shortCode,
			URL:       request.URL,
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
