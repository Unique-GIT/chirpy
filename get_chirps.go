package main

import "net/http"

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
