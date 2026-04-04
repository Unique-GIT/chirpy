package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJsonError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Printf("Sending Error: %v", err)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}

	respondWithJson(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling: %v", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(code)
	w.Write(data)
}
