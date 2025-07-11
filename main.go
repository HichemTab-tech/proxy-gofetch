package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	http.HandleFunc("/fetch", handleFetchProxy)
	log.Println("Image proxy running at http://localhost:8080/fetch")
	log.Fatal(http.ListenAndServe(":3001", nil))
}

func handleFetchProxy(w http.ResponseWriter, r *http.Request) {
	targetURL := r.URL.Query().Get("url")
	if targetURL == "" {
		http.Error(w, "Missing 'url' query parameter", http.StatusBadRequest)
		return
	}

	parsedURL, err := url.Parse(targetURL)
	if err != nil || !parsedURL.IsAbs() {
		http.Error(w, "Invalid target URL", http.StatusBadRequest)
		return
	}

	resp, err := http.Get(targetURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to fetch target", http.StatusBadGateway)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			http.Error(w, "Failed to close response body", http.StatusInternalServerError)
			return
		}
	}(resp.Body)

	// Forward original content type
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))

	// Set CORS header to allow all
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Copy the target stream directly to response
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, "Failed to write target to response", http.StatusInternalServerError)
		return
	}
}
