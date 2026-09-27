package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
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

// sendError sends an error response in JSON format.
func sendError(w http.ResponseWriter, message string, status int) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := struct {
		Error string `json:"error"`
	}{
		Error: message,
	}

	json.NewEncoder(w).Encode(response)
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
			sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Request body structure.
		var request struct {
			URL string `json:"url"`
		}

		// Decode JSON request body.
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			sendError(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Check if URL was provided.
		if request.URL == "" {
			sendError(w, "URL is required", http.StatusBadRequest)
			return
		}
		// Validate the URL.
		parsedURL, err := url.ParseRequestURI(request.URL)
		if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
			sendError(w, "Invalid URL", http.StatusBadRequest)
			return
		}

		// Generate a short code.
		shortCode := generateShortCode()

		// Store the URL using the short code.
		urls[shortCode] = request.URL

		// Create JSON response.
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

	// Redirect endpoint: redirects a short code to the original URL.
	mux.HandleFunc("/short/", func(w http.ResponseWriter, r *http.Request) {

		// Get the short code from the URL.
		shortCode := r.URL.Path[len("/short/"):]

		// Find the original URL.
		originalURL, exists := urls[shortCode]

		// Return 404 if short code does not exist.

		if !exists {
			sendError(w, "Short URL not found", http.StatusNotFound)
			return
		}

		// Redirect to the original URL.
		http.Redirect(w, r, originalURL, http.StatusFound)
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
