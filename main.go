package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/thom151/chirpy/internal/database"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Erorr: %s\n", err)
		return
	}

	dbQueries := database.New(db)

	fmt.Println("Hello World")

	serverMux := http.NewServeMux()

	//config
	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       os.Getenv("PLATFORM"),
		secret:         os.Getenv("SECRET"),
		polkaKey:       os.Getenv("POLKA_KEY"),
	}

	//SERVER
	httpServer := &http.Server{}

	//file HANDLER
	handler := http.FileServer(http.Dir("."))

	serverMux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", handler)))
	serverMux.HandleFunc("GET /api/healthz", handlerReadiness)
	serverMux.HandleFunc("GET /admin/metrics", cfg.handlerCount)
	serverMux.HandleFunc("POST /admin/reset", cfg.handlerReset)
	serverMux.HandleFunc("POST /api/validate_chirp", handlerValidate)
	serverMux.HandleFunc("POST /api/users", cfg.handlerCreateUser)
	serverMux.HandleFunc("POST /api/chirps", cfg.handlerCreateChirp)
	serverMux.HandleFunc("GET /api/chirps", cfg.hanlderGetAllChirps)
	serverMux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handlerGetChirp)
	serverMux.HandleFunc("POST /api/login", cfg.handlerLogin)
	serverMux.HandleFunc("POST /api/refresh", cfg.handlerRefresh)
	serverMux.HandleFunc("POST /api/revoke", cfg.handlerRevoke)
	serverMux.HandleFunc("PUT /api/users", cfg.hanlderUpdateUser)
	serverMux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.handlerDeleteChirp)
	serverMux.HandleFunc("POST /api/polka/webhooks", cfg.handlerUpgradeUser)

	httpServer.Handler = serverMux
	httpServer.Addr = ":8080"

	err = httpServer.ListenAndServe()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
