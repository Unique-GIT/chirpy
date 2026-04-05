package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Unique-GIT/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) validate_chirp(w http.ResponseWriter, r *http.Request) {
	type requestType struct {
		RequestBody string    `json:"body"`
		UserId      uuid.UUID `json:"user_id"`
	}
	type returnType struct {
		Id        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Body      string `json:"body"`
		UserId    string `json:"user_id"`
	}

	// Processing Request
	var request requestType
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithJsonError(w, http.StatusInternalServerError, "ERROR processing Request", err)
		return
	}

	// Check is user exists
	_, err := cfg.db.UsersExists(r.Context(), request.UserId)
	if err != nil {
		respondWithJsonError(w, http.StatusBadRequest, "User Doesn't exist", nil)
		return
	}

	// Validation of Chirp
	// Length check
	if len(request.RequestBody) > 140 {
		// Too long
		respondWithJsonError(w, http.StatusBadRequest, "Chirp Too Long", nil)
		return
	}

	// Check for profane words
	CleanedChirp := validated_body(request.RequestBody)

	// Create this chirp
	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   CleanedChirp,
		UserID: request.UserId,
	})

	respondWithJson(w, http.StatusCreated, returnType{
		Id:        chirp.ID.String(),
		CreatedAt: chirp.CreatedAt.String(),
		UpdatedAt: chirp.UpdatedAt.String(),
		Body:      chirp.Body,
		UserId:    chirp.UserID.String(),
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
