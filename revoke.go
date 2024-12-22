package main

import (
	"net/http"

	"github.com/thom151/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
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

	err = cfg.db.RevokeToken(r.Context(), token.Token)
	if err != nil {
		responseWithError(w, 500, "Token not revoked")
	}

	respondWithJSON(w, 204, nil)
}
