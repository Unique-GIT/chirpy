package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (a *apiConfig) middlewareIncrementMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (a *apiConfig) logHits(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(fmt.Sprintf("Hits: %v", a.fileserverHits.Load())))
	if err != nil {
		log.Printf("Error Writing the Body of the Message in the Metrics Handler: %v", err)
	}
}

func (a *apiConfig) reset(w http.ResponseWriter, r *http.Request) {
	a.fileserverHits.Store(0)
}

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

	config := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	router := http.NewServeMux()
	router.Handle("/app/",
		middlewareLog(
			config.middlewareIncrementMetrics(
				fileHandler,
			),
		),
	)
	router.Handle("/healthz",
		middlewareLog(
			config.middlewareIncrementMetrics(
				http.HandlerFunc(handlerHealth),
			),
		),
	)
	router.Handle("/metrics",
		middlewareLog(
			http.HandlerFunc(config.logHits),
		),
	)
	router.Handle("/reset",
		middlewareLog(
			http.HandlerFunc(config.reset),
		),
	)

	server := http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	server.ListenAndServe()
}
