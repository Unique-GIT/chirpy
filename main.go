package main

import (
	"log"
	"net/http"
)

func handlerHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(http.StatusText(http.StatusOK)))
	if err != nil {
		log.Printf("Error Writing the Body of the Message in the Health Handler: %v", err)
	}

}

func main() {
	const filePathRoot = "."
	const port = "8080"

	router := http.NewServeMux()
	router.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot))))
	router.HandleFunc("/healthz", handlerHealth)

	server := http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	server.ListenAndServe()
}
