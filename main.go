package main

import (
    "crypto/sha1"
    "encoding/hex"
    "fmt"
    "io"
    "log"
    "io/ioutil"
    "net/http"
    "os"
    "path/filepath"
    "flag"
    "github.com/robfig/cron/v3"

    "packhub/modules"
)

var cacheDir *string

// Generates a hash for the URL to use as a cache file name
func urlHash(url string) string {
    h := sha1.New()
    h.Write([]byte(url))
    return hex.EncodeToString(h.Sum(nil))
}

// Check if the response is cached; if so, return the cached data
func getCachedResponse(url string) ([]byte, bool) {
    cacheFile := filepath.Join(*cacheDir, urlHash(url))
    data, err := ioutil.ReadFile(cacheFile)
    if err != nil {
        return nil, false
    }
    return data, true
}

// Cache the response to disk
func cacheResponse(url string, data []byte) {
    cacheFile := filepath.Join(*cacheDir, urlHash(url))
    os.MkdirAll(*cacheDir, os.ModePerm)
    _ = ioutil.WriteFile(cacheFile, data, 0644)
}

// Handles the proxy request
func handleRequest(w http.ResponseWriter, r *http.Request) {
    log.Printf("Received request: %s %s", r.Method, r.URL.Path)
    // Check if we have a cached response
    if data, found := getCachedResponse(r.URL.String()); found {
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
    cacheResponse(r.URL.String(), data)

    // Write response to the client
    w.Write(data)
}

func main() {
    // Define command-line flags
    cacheDir = flag.String("cachedir", "/opt/cache", "Path to cache data")
    port := flag.String("port", "8060", "Port to listen for incomming requests")
    cacheValidTime := flag.String("cachevalidtime", "3600", "Time intervals for deleting older cache - [one day is default value]")
    flag.Parse()

    c := cron.New()
    c.AddFunc("@every 30m", func() {
        cacheCleanup.CacheCleanup(*cacheValidTime,*cacheDir)
    })
    c.Start()
    
    http.HandleFunc("/", handleRequest)
    fmt.Println("Starting proxy server on :", *port)
    log.Fatal(http.ListenAndServe(":" + *port, nil))
}
