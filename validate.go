package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type parameters struct {
	Body string `json:"body"`
}

type ValidateResponse struct {
	Valid bool `json:"valid"`
}

type CleanedResponse struct {
	Cleaned_body string `json:"cleaned_body"`
}

func handlerValidate(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, 500, "Something went wrong")
		return

	}

	if len(params.Body) > 140 {
		responseWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	_ = ValidateResponse{
		Valid: true,
	}

	cleanedBody := findProfane(params.Body)
	respondWithJSON(w, http.StatusOK, CleanedResponse{
		Cleaned_body: cleanedBody,
	})
	return

}

func responseWithError(w http.ResponseWriter, code int, msg string) {
	type ErrorResponse struct {
		Error string `json:"error"`
	}

	decodeError := ErrorResponse{
		Error: msg,
	}
	errDat, _ := json.Marshal(decodeError)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(errDat)
	return

}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	resp, err := json.Marshal(payload)
	if err != nil {
		responseWithError(w, 500, "Something went wrong")
	}
	w.Write([]byte(resp))

}

func findProfane(params string) string {
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(params, " ")
	for i := 0; i < len(words); i++ {
		for _, profaneWord := range profaneWords {
			if strings.ToLower(words[i]) == profaneWord {
				fmt.Println("Profane found")
				words[i] = "****"
			}
		}
	}

	return strings.Join(words, " ")

}
