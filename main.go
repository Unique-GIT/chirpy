package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

func middlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Served Request %s on %s \n", r.Method, r.URL.Path)
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
	router.Handle("GET /api/healthz",
		middlewareLog(
			config.middlewareIncrementMetrics(
				http.HandlerFunc(handlerHealth),
			),
		),
	)
	router.Handle("POST /api/validate_chirp",
		middlewareLog(
			config.middlewareIncrementMetrics(
				http.HandlerFunc(validate_chirp),
			),
		),
	)
	router.Handle("GET /admin/metrics",
		middlewareLog(
			http.HandlerFunc(config.logHits),
		),
	)
	router.Handle("POST /admin/reset",
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
