package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/url"
)

const targetDomain = "https://api.smith.langchain.com"

func handler(w http.ResponseWriter, r *http.Request) {
	// Parse the original request URL
	originalURL, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	slog.Warn(fmt.Sprintf("REQUEST: %+v\n\n", r))

	// Read the request body
	if r.Body != nil {

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(body))
		// Parse the request body as JSON

		if len(body) == 0 {
			slog.Warn("Empty body")
		} else {
			var jsonData map[string]interface{}
			err = json.Unmarshal(body, &jsonData)
			if err != nil {
				slog.Error("Failed to parse request body as JSON")
				return
			}
			prettyJSON, err := json.MarshalIndent(jsonData, "", "  ")
			if err != nil {
				slog.Error("Failed to format JSON data")
				return
			}
			slog.Info(fmt.Sprintf("Request body:\n%s\n", string(prettyJSON)))
		}
	}

	// // Print the formatted JSON data
	// prettyJSON, err := json.MarshalIndent(jsonData, "", "  ")
	// if err != nil {
	// 	http.Error(w, "Failed to format JSON data", http.StatusInternalServerError)
	// 	return
	// }
	// print(fmt.Fprintf(w, "Formatted JSON data:\n%s\n", string(prettyJSON)))

	// Construct the new URL with the target domain
	newURL := targetDomain + originalURL.RequestURI()

	// Create a new request to the target domain
	req, err := http.NewRequest(r.Method, newURL, r.Body)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// Copy headers from the original request to the new request
	req.Header = r.Header

	// Create a new HTTP client
	client := http.Client{}

	// Send the request to the target domain
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to send request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy the response headers to the original response writer
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Set the response status code
	w.WriteHeader(resp.StatusCode)

	// Copy the response body to the original response writer
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, "Failed to copy response", http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/", handler)

	// Start the server on port 8080
	log.Println("Server starting on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
