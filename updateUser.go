package main

import (
	"encoding/json"
	"net/http"

	"github.com/thom151/chirpy/internal/auth"
	"github.com/thom151/chirpy/internal/database"
)

func (cfg *apiConfig) hanlderUpdateUser(w http.ResponseWriter, r *http.Request) {
	bearerKey, err := auth.GetBearerToken(r.Header)
	if err != nil {
		responseWithError(w, 401, "Cannot get bearer")
		return
	}

	jwtID, err := auth.ValidateJWT(bearerKey, cfg.secret)
	if err != nil {
		responseWithError(w, 401, err.Error())
		return
	}

	user, err := cfg.db.GetUserByID(r.Context(), jwtID)
	if err != nil {
		responseWithError(w, 500, "Cannot find user")
		return
	}

	decoder := json.NewDecoder(r.Body)
	email := EmailReq{}
	err = decoder.Decode(&email)
	if err != nil {
		responseWithError(w, 500, "Something went wrong w/ decoding")
		return
	}

	hashed, err := auth.HashPassword(email.Password)
	if err != nil {
		responseWithError(w, 500, "Something went wrong hashing")
		return
	}

	userParams := database.UpdateUserParams{
		HashedPassword: hashed,
		Email:          email.Email,
		ID:             user.ID,
	}
	err = cfg.db.UpdateUser(r.Context(), userParams)
	if err != nil {
		responseWithError(w, 500, "Something went wrong updating the user")
		return
	}

	updatedUser, err := cfg.db.GetUser(r.Context(), email.Email)
	if err != nil {
		responseWithError(w, 500, "Something went wrong getting the updated user")
		return
	}

	userWithToken := UserWithToken{
		ID:          updatedUser.ID,
		CreatedAt:   updatedUser.CreatedAt,
		UpdatedAt:   updatedUser.UpdatedAt,
		Email:       updatedUser.Email,
		Token:       bearerKey,
		IsChirpyRed: updatedUser.IsChirpyRed.Bool,
	}

	respondWithJSON(w, 200, userWithToken)
}
