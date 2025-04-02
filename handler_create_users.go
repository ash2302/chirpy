package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) createUsersHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		log.Printf("Error decoding body: %s", err)
		respondWithError(w, http.StatusBadRequest, "Bad Request")
		return
	}
	if params.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Email is required")
		return
	}

	user, err := cfg.dbQueries.CreateUser(r.Context(), params.Email)
	if err != nil {
		log.Printf("Error creating user with email %s: %s", params.Email, err)
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	responseData := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	respondWithJSON(w, http.StatusCreated, responseData)
}
