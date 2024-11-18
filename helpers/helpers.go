package helpers

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/robfig/cron/v3"
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
func CacheResponse(url, cacheDir string, data []byte) {
	cacheFile := filepath.Join(cacheDir, UrlHash(url))
	os.MkdirAll(cacheDir, os.ModePerm)
	_ = os.WriteFile(cacheFile, data, 0644)
}

func cacheCleanup(cacheValidTime, cacheDir string) {

	// Get the current time
	now := time.Now()

	// Read the contents of the directory
	files, err := os.ReadDir(cacheDir)
	if err != nil {
		log.Fatal(err)
	}

	// Print the list of files
	for _, file := range files {
		filePath := cacheDir + "/" + file.Name()

		// Get file information
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			log.Fatal(err)
		}

		cachevalidtime, _ := strconv.Atoi(cacheValidTime)

		// Check if the file was modifed more than one minute
		if now.Sub(fileInfo.ModTime()) > time.Duration(cachevalidtime)*time.Minute {
			os.Remove(filePath)
			fmt.Println("file: ", file.Name(), " ***DELETED*** ")
		}
	}
}

func CacheCronJob(expiration_time, cacheDir string) {
	c := cron.New()
	c.AddFunc("@every 30m", func() {
		cacheCleanup(expiration_time, cacheDir)
	})
	c.Start()
}
