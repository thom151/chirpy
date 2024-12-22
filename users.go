package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/thom151/chirpy/internal/auth"
	"github.com/thom151/chirpy/internal/database"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

type EmailReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserWithToken struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("calling")
	decoder := json.NewDecoder(r.Body)
	email := EmailReq{}
	err := decoder.Decode(&email)
	if err != nil {
		responseWithError(w, 500, "Something went wrong w/ decoding")
		return
	}

	dbUser, err := cfg.db.CreateUser(r.Context(), email.Email)
	if err != nil {
		responseWithError(w, 500, "Something went wrong w/ creating user")
		fmt.Printf("Error: %v\n", err)
		return
	}

	hashed, err := auth.HashPassword(email.Password)
	if err != nil {
		responseWithError(w, 500, "Something went wrong w/ hasing password.")
		return
	}

	passParams := database.SetPasswordParams{
		HashedPassword: hashed,
		ID:             dbUser.ID,
	}
	err = cfg.db.SetPassword(r.Context(), passParams)
	if err != nil {
		responseWithError(w, 500, "Something went wrong with setting password.")
	}

	fmt.Printf("Password has been set for %s\n", dbUser.Email)

	myUser := User{
		ID:          dbUser.ID,
		CreatedAt:   dbUser.CreatedAt,
		UpdatedAt:   dbUser.UpdatedAt,
		Email:       dbUser.Email,
		IsChirpyRed: dbUser.IsChirpyRed.Bool,
	}

	respondWithJSON(w, 201, myUser)

}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	email := EmailReq{}
	err := decoder.Decode(&email)
	if err != nil {
		responseWithError(w, 500, "Error decoding request body")
		return
	}

	dbUser, err := cfg.db.GetUser(r.Context(), email.Email)
	if err != nil {
		responseWithError(w, 500, "User Does not exiist i reckon!")
		return
	}

	err = auth.CheckPasswordHash(email.Password, dbUser.HashedPassword)
	if err != nil {
		responseWithError(w, 401, "Unauthorized")
		return
	}

	userJWT, err := auth.MakeJWT(dbUser.ID, cfg.secret, time.Hour)
	if err != nil {
		responseWithError(w, 500, "Something went wrong with getting token")
		return
	}

	myUser := User{
		ID:          dbUser.ID,
		CreatedAt:   dbUser.CreatedAt,
		UpdatedAt:   dbUser.UpdatedAt,
		Email:       dbUser.Email,
		IsChirpyRed: dbUser.IsChirpyRed.Bool,
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		responseWithError(w, 500, "Something went wrong w/ making refresh token")
		return
	}

	refreshParams := database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    dbUser.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
		RevokedAt: sql.NullTime{},
	}

	_, err = cfg.db.CreateRefreshToken(r.Context(), refreshParams)
	if err != nil {
		responseWithError(w, 500, "Error creating ref token")
		return
	}
	userWithToken := UserWithToken{
		ID:           myUser.ID,
		CreatedAt:    myUser.CreatedAt,
		UpdatedAt:    myUser.UpdatedAt,
		Email:        myUser.Email,
		Token:        userJWT,
		RefreshToken: refreshToken,
		IsChirpyRed:  myUser.IsChirpyRed,
	}

	respondWithJSON(w, 200, userWithToken)

}
