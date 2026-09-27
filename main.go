package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

//go:embed frontend
var frontendFiles embed.FS

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

	// Stores each IP's request count and the start of its 1-minute window.
	type rateLimit struct {
		count       int
		windowStart time.Time
	}
	rateLimits := make(map[string]rateLimit)

	// Protects both maps because HTTP handlers can run at the same time.
	var mu sync.Mutex

	// Create a new router.
	mux := http.NewServeMux()

	frontend, err := fs.Sub(frontendFiles, "frontend")
	if err != nil {
		panic(err)
	}
	mux.Handle("/", http.FileServer(http.FS(frontend)))

	// Health endpoint.
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	// Shorten endpoint.
	mux.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {

		// Only allow POST requests.
		if r.Method != http.MethodPost {
			sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Render forwards the visitor's IP in X-Forwarded-For.
		ip := r.Header.Get("X-Forwarded-For")
		if ip != "" {
			// The header can contain multiple IPs; the first is the visitor.
			ip = strings.TrimSpace(strings.Split(ip, ",")[0])
		} else {
			// Use the direct connection IP when no proxy header is present.
			ip, _, err = net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			}
		}

		// Start a new 1-minute window when this IP makes its first request
		// or its previous window has expired.
		now := time.Now()
		mu.Lock()
		limit, exists := rateLimits[ip]
		if !exists || now.Sub(limit.windowStart) >= time.Minute {
			limit = rateLimit{count: 1, windowStart: now}
		} else if limit.count >= 5 {
			mu.Unlock()
			sendError(w, "Rate limit exceeded. Try again later.", http.StatusTooManyRequests)
			return
		} else {
			limit.count++
		}
		rateLimits[ip] = limit
		mu.Unlock()

		// Request body structure.
		var request struct {
			URL string `json:"url"`
		}

		// Decode JSON request body.
		err = json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			sendError(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Check if URL was provided.
		if request.URL == "" {
			sendError(w, "URL is required", http.StatusBadRequest)
			return
		}

		// Validate URL.
		parsedURL, err := url.ParseRequestURI(request.URL)
		if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
			sendError(w, "Invalid URL", http.StatusBadRequest)
			return
		}

		// Generate a unique short code.
		mu.Lock()
		shortCode := generateShortCode()

		for {
			_, exists := urls[shortCode]

			if !exists {
				break
			}

			shortCode = generateShortCode()
		}

		// Store the URL.
		urls[shortCode] = request.URL
		mu.Unlock()

		// Create response.
		response := struct {
			Message   string `json:"message"`
			ShortCode string `json:"short_code"`
			URL       string `json:"url"`
		}{
			Message:   "URL shortened successfully",
			ShortCode: shortCode,
			URL:       request.URL,
		}

		// Send JSON response.
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Redirect endpoint.
	mux.HandleFunc("/short/", func(w http.ResponseWriter, r *http.Request) {

		// Get short code from URL.
		shortCode := r.URL.Path[len("/short/"):]

		// Find original URL.
		mu.Lock()
		originalURL, exists := urls[shortCode]
		mu.Unlock()

		if !exists {
			sendError(w, "Short URL not found", http.StatusNotFound)
			return
		}

		// Redirect to original URL.
		http.Redirect(w, r, originalURL, http.StatusFound)
	})

	// Render provides the port in the PORT environment variable.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Bind on all network interfaces so Render can forward public traffic.
	address := "0.0.0.0:" + port
	fmt.Println("Server running on port", port)

	err = http.ListenAndServe(address, mux)

	if err != nil {
		fmt.Println("Server error:", err)
	}
}

func main() {

	// Start the server.
	startServer()
}
