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

func middlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Served Request %s on %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func main() {
	const filePathRoot = "."
	const port = "8080"

	fileHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot)))
	router := http.NewServeMux()
	router.Handle("/app/", middlewareLog(fileHandler))
	router.Handle("/healthz", middlewareLog(http.HandlerFunc(handlerHealth)))

	server := http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	server.ListenAndServe()
}
