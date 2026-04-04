package main

import "net/http"

func (a *apiConfig) reset(w http.ResponseWriter, r *http.Request) {
	a.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
