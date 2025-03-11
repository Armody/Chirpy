package main

import (
	"net/http"
	"time"

	"github.com/Armody/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, req *http.Request) {
	type returnVals struct {
		Token string `json:"token"`
	}
	refString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get refresh token", err)
		return
	}

	refToken, err := cfg.db.GetRefreshToken(req.Context(), refString)
	if err != nil || refToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
		return
	}

	newToken, err := auth.MakeJWT(refToken.UserID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create new token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		Token: newToken,
	})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, req *http.Request) {
	refString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get refresh token", err)
		return
	}

	if err := cfg.db.RevokeRefreshToken(req.Context(), refString); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke refresh token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
