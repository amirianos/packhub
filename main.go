package main

import (
	"flag"
	"fmt"
	"net/http"
	"packhub/handlers"
	"packhub/helpers"
)

func main() {
	// Define command-line flags
	cacheDir := flag.String("cachedir", "/opt/cache", "Path to cache data")
	port := flag.String("port", "8060", "Port to listen for incomming requests")
	cacheValidTime := flag.String("cachevalidtime", "3600", "Time intervals for deleting older cache - [one day is default value]")
	flag.Parse()
	helpers.CacheCronJob(*cacheValidTime, *cacheDir)

	handlers := handlers.New(*cacheDir)
	server := http.Server{
		Addr:    fmt.Sprintf(": %s", *port),
		Handler: getRoutes(handlers),
	}

	server.ListenAndServe()
}
