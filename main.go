package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const (
	pubDevURL    = "https://pub.dev"
	cacheDir     = "./cache" // Directory to store cached packages
	port         = ":8060"
)

func main() {
	http.HandleFunc("/", handleRequest)
	fmt.Println("Proxy server started on port", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Println("Server error:", err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	packagePath := r.URL.Path
	cachedFilePath := filepath.Join(cacheDir, packagePath)

	if _, err := os.Stat(cachedFilePath); err == nil {
		// Serve from cache if it exists
		http.ServeFile(w, r, cachedFilePath)
		fmt.Println("Served from cache:", packagePath)
		return
	}

	// Fetch from pub.dev if not in cache
	fmt.Println("Fetching from pub.dev:", packagePath)
	resp, err := http.Get(pubDevURL + packagePath)
	if err != nil {
		http.Error(w, "Unable to fetch package", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Create cache directories if needed
	if err := os.MkdirAll(filepath.Dir(cachedFilePath), os.ModePerm); err != nil {
		http.Error(w, "Unable to create cache directories", http.StatusInternalServerError)
		return
	}

	// Cache the package locally
	file, err := os.Create(cachedFilePath)
	if err != nil {
		http.Error(w, "Unable to save package", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		http.Error(w, "Error saving package", http.StatusInternalServerError)
		return
	}

	// Serve the fetched package to the client
	http.ServeFile(w, r, cachedFilePath)
}

