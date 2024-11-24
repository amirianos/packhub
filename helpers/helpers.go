package helpers

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"packhub/models"
	"path/filepath"
	"time"

	"github.com/robfig/cron/v3"
	"gopkg.in/yaml.v3"
)

// Generates a hash for the URL to use as a cache file name
func UrlHash(url string) string {
	h := sha1.New()
	h.Write([]byte(url))
	return hex.EncodeToString(h.Sum(nil))
}

// Check if the response is cached; if so, return the cached data
func GetCachedResponse(url, cacheDir string) ([]byte, bool) {
	cacheFile := filepath.Join(cacheDir, UrlHash(url))
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil, false
	}
	return data, true
}

// Cache the response to disk
func CacheResponse(url, cacheDir string, data []byte) error {
	cacheFile := filepath.Join(cacheDir, UrlHash(url))
	err := os.MkdirAll(cacheDir, os.ModePerm)
	if err != nil {
		return err
	}
	return os.WriteFile(cacheFile, data, 0644)

}

func cacheCleanup(cacheValidTime int, cacheDir string) {
	now := time.Now()
	files, err := os.ReadDir(cacheDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		filePath := cacheDir + "/" + file.Name()
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			log.Fatal(err)
		}
		if now.Sub(fileInfo.ModTime()) > time.Duration(cacheValidTime)*time.Hour {
			os.Remove(filePath)
			fmt.Println("file: ", file.Name(), " ***DELETED*** ")
		}
	}
}

func CacheCronJob(remoteRepos map[string]*models.RemoteRepository) {
	c := cron.New()
	for _, remoteRepo := range remoteRepos {
		c.AddFunc(fmt.Sprintf("@every %s", remoteRepo.CacheExpirationCheck), func() {
			cacheCleanup(remoteRepo.CacheValidTime, remoteRepo.CacheDirectory)
		})
	}
	c.Start()
}

func ParseYaml(path string) (map[string]*models.RemoteRepository, error) {

	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	repos := make(map[string]*models.RemoteRepository)

	err = yaml.Unmarshal(yamlFile, repos)
	if err != nil {
		return nil, err
	}
	return repos, nil
}

func MessageGenerator(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)
	msg := struct {
		Message string `json:"message"`
	}{
		Message: message,
	}
	jsonResponse, _ := json.Marshal(msg)
	w.Write(jsonResponse)
}
