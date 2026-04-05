package main

import (
	"encoding/json"
	"net/http"

	"github.com/Unique-GIT/chirpy/internal/auth"
	"github.com/Unique-GIT/chirpy/internal/database"
)

func (a *apiConfig) handlerUser(w http.ResponseWriter, r *http.Request) {
	type requestType struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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

	// Hash Password
	hash, err := auth.HashPassword(request.Password)
	if err != nil {
		respondWithJsonError(w, http.StatusInternalServerError, "Error Hashing Password", err)
		return
	}

	// Process Request
	user, err := a.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          request.Email,
		HashedPassword: hash,
	})
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

func (a *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type requestType struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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

	// Find user
	user, err := a.db.GetUserByEmail(r.Context(), request.Email)
	if err != nil {
		respondWithJsonError(w, http.StatusUnauthorized, "User Does Not Exist", nil)
		return
	}

	// Compare password
	correctPassword, err := auth.CheckPasswordHash(request.Password, user.HashedPassword)
	if err != nil {
		respondWithJsonError(w, http.StatusInternalServerError, "Error comparing hashes of password", err)
		return
	}
	if !correctPassword {
		respondWithJsonError(w, http.StatusUnauthorized, "Wrong Password", nil)
		return
	}

	respondWithJson(w, http.StatusOK, responseType{
		Id:        user.ID.String(),
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
		Email:     user.Email,
	})
}
