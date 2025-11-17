package main

import (
	"fmt"
	"net/http"
)

// This function handles requests to the root URL ("/")
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests to "/"
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Send a simple text response
	fmt.Fprintf(w, "Hello, synced music queue! Server is running.")
}

func main() {
	// Register the handler for the "/" route
	http.HandleFunc("/", homeHandler)

	// Start the server on port 8080
	fmt.Println("Server starting on http://localhost:8080")
	fmt.Println("Open your browser and go to that address!")

	// ListenAndServe starts the HTTP server
	// If there's an error (e.g. port in use), it will print it
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Failed to start server:", err)
	}
}