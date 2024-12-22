package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/thom151/chirpy/internal/auth"
	"github.com/thom151/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

type ChirpReq struct {
	Body string `json:"body"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {

	bearerKey, err := auth.GetBearerToken(r.Header)
	if err != nil {
		responseWithError(w, 401, "Cannot get bearer")
		return
	}

	fmt.Printf("bearerere: %s\n", bearerKey)
	jwtID, err := auth.ValidateJWT(bearerKey, cfg.secret)
	if err != nil {
		responseWithError(w, 401, err.Error())
		return
	}

	decoder := json.NewDecoder(r.Body)
	chirp := ChirpReq{}
	err = decoder.Decode(&chirp)
	if err != nil {
		responseWithError(w, 500, "Something went wrong with decoding chirp")
		return
	}

	if len(chirp.Body) > 140 {
		responseWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanedChirp := findProfane(chirp.Body)

	dbChirpParams := database.CreateChirpParams{
		Body:   cleanedChirp,
		UserID: jwtID,
	}

	dbChirp, err := cfg.db.CreateChirp(r.Context(), dbChirpParams)
	if err != nil {
		responseWithError(w, 500, "Internal server error: Something went wrong with creating chirp")
		return
	}

	myChirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}

	respondWithJSON(w, 201, myChirp)

}

func respondMyChirps(chirps []database.Chirp) []Chirp {
	myChirps := make([]Chirp, 0)
	for _, chirp := range chirps {

		myChirp := Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
		myChirps = append(myChirps, myChirp)
	}

	return myChirps

}

func (cfg *apiConfig) hanlderGetAllChirps(w http.ResponseWriter, r *http.Request) {
	myChirps := make([]Chirp, 0)

	authorId := r.URL.Query().Get("author_id")
	order := r.URL.Query().Get("sort")

	if authorId != "" {
		userUUID, err := uuid.Parse(authorId)
		if err != nil {
			responseWithError(w, 400, "Error parsing id")
			return
		}

		_, err = cfg.db.GetUserByID(r.Context(), userUUID)
		if err != nil {
			responseWithError(w, 400, "Cannot find user")
			return
		}

		if order == "desc" {
			chirps, err := cfg.db.GetAllChirpsFromUserDesc(r.Context(), userUUID)
			if err != nil {
				responseWithError(w, 500, "Something went wrong getting chirps for users")
				return
			}

			myChirps = respondMyChirps(chirps)

			respondWithJSON(w, 200, myChirps)
			return

		}

		chirps, err := cfg.db.GetAllChirpsFromUser(r.Context(), userUUID)
		if err != nil {
			responseWithError(w, 500, "Something went wrong getting chirps for users")
			return
		}

		myChirps = respondMyChirps(chirps)

		respondWithJSON(w, 200, myChirps)
		return

	}

	if order == "desc" {
		chirps, err := cfg.db.GetAllChirpsDesc(r.Context())
		if err != nil {
			responseWithError(w, 500, "Something went wrong with getting chirps")
			return
		}

		myChirps = respondMyChirps(chirps)
		respondWithJSON(w, 200, myChirps)
		return
	}

	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		responseWithError(w, 500, "Something went wrong with getting chirps")
		return
	}

	myChirps = respondMyChirps(chirps)
	respondWithJSON(w, 200, myChirps)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	fmt.Printf("Chirp Id: %s\n", chirpID)

	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		responseWithError(w, 500, "Invalid chirp id")
		return
	}

	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpUUID)
	if err != nil {
		responseWithError(w, 404, "Cannot find chirp")
		return
	}

	myChirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}

	respondWithJSON(w, 200, myChirp)

}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {

	bearerKey, err := auth.GetBearerToken(r.Header)
	if err != nil {
		responseWithError(w, 401, "Something went wrong getting bearer of delete Chirp")
		return
	}

	jwtID, err := auth.ValidateJWT(bearerKey, cfg.secret)
	if err != nil {
		respondWithJSON(w, 403, "Unauthorized access")
		return
	}

	chirpID := r.PathValue("chirpID")

	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		responseWithError(w, 500, "Something went wrong parsing chirp id")
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpUUID)
	if err != nil {
		responseWithError(w, 500, "Something went wrong getting chirp from del")
		return
	}

	if chirp.UserID != jwtID {
		responseWithError(w, 403, "Unauthorized acces, not creator")
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), chirpUUID)
	if err != nil {
		responseWithError(w, 404, "Something went wrong deleting chirp")
		return
	}

	w.WriteHeader(204)

}
