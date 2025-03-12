package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Armody/Chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerWebhookUpgrade(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserId uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	apiKey, err := auth.GetApiKey(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't extract API key", err)
		return
	}

	if apiKey != cfg.polkaApi {
		respondWithError(w, http.StatusUnauthorized, "API key is invalid", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if _, err := cfg.db.UpgradeToChirpyRed(req.Context(), params.Data.UserId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "User not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error upgrading user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
