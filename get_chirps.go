package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithJsonError(w, http.StatusInternalServerError, "Error Getting Chirps", err)
		return
	}

	var results []Chirp
	for _, chirp := range chirps {
		results = append(results, Chirp{
			Id:        chirp.ID.String(),
			CreatedAt: chirp.CreatedAt.String(),
			UpdatedAt: chirp.UpdatedAt.String(),
			Body:      chirp.Body,
			UserId:    chirp.UserID.String(),
		})
	}

	respondWithJson(w, http.StatusOK, results)
}

func (cfg *apiConfig) getChirpById(w http.ResponseWriter, r *http.Request) {
	stringId := r.PathValue("ChirpsId")

	id, err := uuid.Parse(stringId)
	if err != nil {
		respondWithJsonError(w, http.StatusBadRequest, "Id Not correct", err)
		return
	}

	chirp, err := cfg.db.GetChirpById(r.Context(), id)
	if err != nil {
		respondWithJsonError(w, http.StatusNotFound, "Failed To get Chirp", err)
		return
	}

	respondWithJson(w, http.StatusOK, Chirp{
		Id:        chirp.ID.String(),
		CreatedAt: chirp.CreatedAt.String(),
		UpdatedAt: chirp.UpdatedAt.String(),
		Body:      chirp.Body,
		UserId:    chirp.UserID.String(),
	})
}
