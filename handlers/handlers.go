package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"packhub/helpers"
	"packhub/models"
	"strings"

	"github.com/go-chi/chi/v5"
)

type Repository struct {
	remoteRepositories map[string]*models.RemoteRepository
}

func New(remoteRepos map[string]*models.RemoteRepository) *Repository {
	return &Repository{
		remoteRepositories: remoteRepos,
	}
}

func (repo *Repository) Home(w http.ResponseWriter, r *http.Request) {
	requestedRepo := chi.URLParam(r, "provider")
	provider, ok := repo.remoteRepositories[requestedRepo]
	if !ok {
		helpers.MessageGenerator(w, "The requested provider is not supported", http.StatusBadRequest)
		return
	}

	providerRemoteAddress := provider.Address
	requestedPackage := strings.Split(r.URL.Path, requestedRepo)[1]
	cacheDir := provider.CacheDirectory
	if data, found := helpers.GetCachedResponse(requestedPackage, cacheDir); found {
		log.Printf("%s%s is already cached", providerRemoteAddress, requestedPackage)
		w.Write(data)
		return
	}
	// Forward request to the target server
	resp, err := http.Get(providerRemoteAddress + requestedPackage)
	if err != nil {
		log.Println(err.Error())
		helpers.MessageGenerator(w, fmt.Sprintf("Failed to fetch %s%s", providerRemoteAddress, requestedPackage), http.StatusBadRequest)
		return
	}
	if resp.StatusCode != 200 {
		helpers.MessageGenerator(w, fmt.Sprintf("Failed to fetch %s%s", providerRemoteAddress, requestedPackage), resp.StatusCode)
		return
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	err = helpers.CacheResponse(requestedPackage, cacheDir, data)
	if err != nil {
		log.Println(err.Error())
	}
	w.Write(data)
}
