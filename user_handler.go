package main

import (
	"encoding/json"
	"net/http"
)

func (a *apiConfig) handlerUser(w http.ResponseWriter, r *http.Request) {
	type requestType struct {
		Email string `json:"email"`
	}
	type responseType struct {
		Id        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Email     string `json:"email"`
	}

	// Get Request
	var request requestType
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithJsonError(w, http.StatusBadRequest, "Wrong Input Format", err)
		return
	}

	// Process Request
	user, err := a.db.CreateUser(r.Context(), request.Email)
	if err != nil {
		respondWithJsonError(w, http.StatusInternalServerError, "Internal Server Error", err)
		return
	}

	// Respond
	respondWithJson(w, http.StatusCreated, responseType{
		Id:        user.ID.String(),
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
		Email:     user.Email,
	})
}
