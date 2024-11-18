package handlers

import (
	"io"
	"log"
	"net/http"
	"packhub/helpers"
)

type Repository struct {
	cacheDir string
}

func New(cacheDirectory string) *Repository {
	return &Repository{
		cacheDir: cacheDirectory,
	}
}

func (repo *Repository) Home(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %s %s", r.Method, r.URL.Path)
	// Check if we have a cached response
	if data, found := helpers.GetCachedResponse(r.URL.String(), repo.cacheDir); found {
		w.Write(data)
		return
	}
	// Forward request to the target server
	resp, err := http.Get("https://pub.dev" + r.URL.Path)
	if err != nil {
		http.Error(w, "Failed to fetch", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	// Read and cache the response
	data, _ := io.ReadAll(resp.Body)
	helpers.CacheResponse(r.URL.String(), repo.cacheDir, data)
	// Write response to the client
	w.Write(data)
}
