package handlers

import (
	"encoding/json"
	"fmt"
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
	// log.Printf("Received request: %s %s\n", r.Method, r.URL.Path)
	// log.Printf("Path to the corresponding package is: https://pub.dev%s", r.URL.Path)
	// Check if we have a cached response
	if data, found := helpers.GetCachedResponse(r.URL.String(), repo.cacheDir); found {
		log.Printf("https://pub.dev%s is already cached" , r.URL.Path)
		w.Write(data)
		return
	}
	// Forward request to the target server
	resp, err := http.Get("https://pub.dev" + r.URL.Path)
	if err != nil {
		log.Println(err.Error())
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		msg := struct {
			Message string `json:"message"`
		}{
			Message: fmt.Sprintf("Failed to fetch https://pub.dev%s", r.URL.Path),
		}
		jsonResponse, _ := json.Marshal(msg)
		w.Write(jsonResponse)
		return
	}
	defer resp.Body.Close()
	// Read and cache the response
	data, _ := io.ReadAll(resp.Body)
	err = helpers.CacheResponse(r.URL.String(), repo.cacheDir, data)
	if err != nil {
		log.Println(err.Error())
	}
	// Write response to the client
	w.Write(data)
}
