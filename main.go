package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

var allowedDomains []string
var allowedRegexPatterns []*regexp.Regexp

func init() {
	// Get domains from environment variable
	domainsEnv := os.Getenv("ALLOWED_DOMAINS")
	if domainsEnv == "" {
		// If ALLOWED_DOMAINS is set to "", disallow all domains
		allowedDomains = []string{"NO_DOMAIN_ALLOWED"}
		return
	} else if domainsEnv == "*" {
		// If ALLOWED_DOMAINS is set to "*", allow all domains
		allowedDomains = []string{"*"}
		return
	}

	// Split domains by comma and process each
	domains := strings.Split(domainsEnv, ",")
	for _, domain := range domains {
		domain = strings.TrimSpace(domain)
		if strings.Contains(domain, "*") {
			// Convert wildcard domain to regex pattern
			pattern := strings.Replace(domain, ".", "\\.", -1)
			pattern = strings.Replace(pattern, "*", "[a-zA-Z0-9-]+", -1)
			pattern = "^https?://" + pattern + "$"
			if regexPattern, err := regexp.Compile(pattern); err == nil {
				allowedRegexPatterns = append(allowedRegexPatterns, regexPattern)
			}
		} else {
			allowedDomains = append(allowedDomains, domain)
		}
	}
}

func isAllowedOrigin(origin string) bool {
	if len(allowedDomains) == 1 && allowedDomains[0] == "*" {
		return true
	}

	// Check exact matches
	for _, domain := range allowedDomains {
		if origin == domain {
			return true
		}
	}

	// Check regex patterns
	for _, pattern := range allowedRegexPatterns {
		if pattern.MatchString(origin) {
			return true
		}
	}

	return false
}

func handleCORS(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin != "" && isAllowedOrigin(origin) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "3600")
	}
}

func main() {
	http.HandleFunc("/fetch", handleFetchProxy)
	log.Println("Image proxy running at http://localhost:8080/fetch")
	log.Fatal(http.ListenAndServe(":3001", nil))
}

func handleFetchProxy(w http.ResponseWriter, r *http.Request) {
	// Handle CORS preflight request
	if r.Method == "OPTIONS" {
		handleCORS(w, r)
		return
	}

	// Handle CORS for actual request
	handleCORS(w, r)

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

	// Copy the target stream directly to response
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, "Failed to write target to response", http.StatusInternalServerError)
		return
	}
}
