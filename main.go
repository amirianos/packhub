package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"packhub/handlers"
	"packhub/helpers"
)

func main() {
	port := flag.String("port", "8060", "Port to listen for incomming requests")
	flag.Parse()
	remoteRepos, err := helpers.ParseYaml("repositories.yml")
	if err != nil {
		log.Fatalln("Could not parse remote repositories from a yaml file")
	}
	helpers.CacheCronJob(remoteRepos)
	handlers := handlers.New(remoteRepos)
	server := http.Server{
		Addr:    fmt.Sprintf(": %s", *port),
		Handler: getRoutes(handlers),
	}

	log.Println("proxy server is running on port :", *port)
	server.ListenAndServe()
}
