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

func (a *apiConfig) logHits(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(fmt.Sprintf("Hits: %v", a.fileserverHits.Load())))
	if err != nil {
		log.Printf("Error Writing the Body of the Message in the Metrics Handler: %v \n ", err)
	}

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func (a *apiConfig) middlewareIncrementMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
