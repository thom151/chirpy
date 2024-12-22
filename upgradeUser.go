package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/thom151/chirpy/internal/auth"
)

type UpgradeReq struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) handlerUpgradeUser(w http.ResponseWriter, r *http.Request) {

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		responseWithError(w, 401, "Unauthorized")
		return
	}

	if apiKey != cfg.polkaKey {
		responseWithError(w, 401, "Unauthorized")
		return
	}

	decoder := json.NewDecoder(r.Body)
	upgradeDat := UpgradeReq{}

	err = decoder.Decode(&upgradeDat)
	if err != nil {
		responseWithError(w, 500, "Something went wrong decoding upgrade request")
		return
	}

	if upgradeDat.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	userUUID, err := uuid.Parse(upgradeDat.Data.UserID)
	if err != nil {
		responseWithError(w, 500, "Something went wroing with parsing id")
		return
	}

	_, err = cfg.db.GetUserByID(r.Context(), userUUID)
	if err != nil {
		responseWithError(w, 404, "Cannot find user")
		return
	}

	err = cfg.db.UpgradeUser(r.Context(), userUUID)
	if err != nil {
		responseWithError(w, 404, "Something went wrong upgrading user"+err.Error())
		return
	}

	w.WriteHeader(204)
	return
}
