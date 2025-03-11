package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Armody/Chirpy/internal/auth"
	"github.com/Armody/Chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type returnVals struct {
		User
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hashedPwd, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
	}

	dbUser, err := cfg.db.CreateUser(req.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPwd,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, returnVals{
		User: User{
			ID:        dbUser.ID,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			Email:     dbUser.Email,
		},
	})
}

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	AccToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	userID, err := auth.ValidateJWT(AccToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token", err)
		return
	}

	if params.Email != "" {
		if err := cfg.db.UpdateEmail(req.Context(), database.UpdateEmailParams{
			ID:    userID,
			Email: params.Email,
		}); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't update Email", err)
			return
		}
	}

	if params.Password != "" {
		newPassword, err := auth.HashPassword(params.Password)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
			return
		}

		if err := cfg.db.UpdatePassword(req.Context(), database.UpdatePasswordParams{
			ID:             userID,
			HashedPassword: newPassword,
		}); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't update password", err)
			return
		}
	}

	dbUser, err := cfg.db.GetUserById(req.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "No user found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:        dbUser.ID,
		Email:     dbUser.Email,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	})
}
