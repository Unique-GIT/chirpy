package main

import (
	"encoding/json"
	"net/http"
)

func validate_chirp(w http.ResponseWriter, r *http.Request) {
	type requestType struct {
		RequestBody string `json:"body"`
	}
	type returnType struct {
		Validity bool `json:"valid"`
	}

	// Processing Request
	var request requestType
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithJsonError(w, http.StatusInternalServerError, "ERROR processing Request", err)
		return
	}

	// Length check
	if len(request.RequestBody) > 140 {
		// Too long
		respondWithJsonError(w, http.StatusBadRequest, "Chirp Too Long", nil)
		return
	}

	respondWithJson(w, http.StatusOK, returnType{
		Validity: true,
	})
}
