package main

import (
	"log"
	"net/http"
)

func main() {
	const filePathRoot = "."
	const port = "8080"

	router := http.NewServeMux()
	router.Handle("/", http.FileServer(http.Dir(filePathRoot)))

	server := http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	server.ListenAndServe()
}
