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

func (a *apiConfig) handlerDeleteUser(w http.ResponseWriter, r *http.Request) {
	if a.platform != "dev" {
		respondWithJsonError(w, http.StatusForbidden, "Deletion Not Allowed", nil)
		return
	}

	a.fileserverHits.Store(0)
	err := a.db.DeleteUsers(r.Context())
	if err != nil {
		respondWithJsonError(w, http.StatusInternalServerError, "Error Deleting users", err)
		return
	}

	type responseType struct {
		Message string `json:"msg"`
	}
	respondWithJson(w, http.StatusOK, responseType{
		Message: "Deleted All users",
	})
}
