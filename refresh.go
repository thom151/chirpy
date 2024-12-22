package main

import (
	"net/http"
	"time"

	"github.com/thom151/chirpy/internal/auth"
)

type TokenResp struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		responseWithError(w, 500, "Something went wrong with getting bearer")
		return
	}

	token, err := cfg.db.GetToken(r.Context(), refreshToken)
	if err != nil {
		responseWithError(w, 401, "Token does not exist")
		return
	}
	if token.ExpiresAt.Before(time.Now()) {
		responseWithError(w, 401, "Token expired")
		return
	}

	newJWT, err := auth.MakeJWT(token.UserID, cfg.secret, time.Hour)
	if err != nil {
		responseWithError(w, 500, "Error creating new access token")
		return
	}

	if token.RevokedAt.Valid { // or however your nullable timestamp is structured
		responseWithError(w, 401, "Token has been revoked")
		return
	}

	tokenResp := TokenResp{
		Token: newJWT, // Return the new JWT, not the refresh token
	}
	respondWithJSON(w, 200, tokenResp)
}
