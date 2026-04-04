package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func validate_chirp(w http.ResponseWriter, r *http.Request) {
	type requestType struct {
		RequestBody string `json:"body"`
	}
	type returnType struct {
		Response string `json:"cleaned_body"`
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

	// Check for profane words
	respondWithJson(w, http.StatusOK, returnType{
		Response: validated_body(request.RequestBody),
	})
}

func validated_body(input string) string {
	splitBody := strings.Split(input, " ")
	badWords := []string{"kerfuffle", "sharbert", "fornax"}

	newStrings := []string{}
	for _, word := range splitBody {
		isBad := false
		for _, badWord := range badWords {
			if strings.ToLower(word) == badWord {
				isBad = true
			}
		}

		if isBad {
			newStrings = append(newStrings, "****")
		} else {
			newStrings = append(newStrings, word)
		}
	}

	cleaned := strings.Join(newStrings, " ")
	return cleaned
}
