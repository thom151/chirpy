package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/thom151/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	secret         string
	polkaKey       string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	countStr := fmt.Sprintf(`<html>
								<body>
								    <h1>Welcome, Chirpy Admin</h1>
									<p>Chirpy has been visited %d times!</p>
								</body>
							</html>`,
		cfg.fileserverHits.Load())
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(countStr))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {

	if cfg.platform != "dev" {
		responseWithError(w, 403, "Forbidden")
		return
	}
	cfg.fileserverHits.Store(0)
	err := cfg.db.DeleteUsers(r.Context())
	if err != nil {
		log.Fatalf("Error deleting users: %v\n", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0 and database reset to initial state."))

}
