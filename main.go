package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"saurabh/chirpy.com/m/internal/database"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	secret := os.Getenv("AUTH_SECRET")

	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		return
	}

	dbQueries := database.New(db)

	apiCfg := apiConfig{
		fileServerHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
		secret:         secret,
	}

	const port = "8080"
	mux := http.NewServeMux()

	httpServer := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}
	appHandler := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", apiCfg.MiddlewareMetricsIncrement(appHandler))
	mux.HandleFunc("GET /admin/metrics", apiCfg.MetricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.ResetMetricsHandler)

	mux.HandleFunc("GET /api/healthz", HealthCheckHandler)
	mux.HandleFunc("POST /api/chirps", apiCfg.CreateChirpHandler)
	mux.HandleFunc("POST /api/users", apiCfg.CreateUserHandler)
	mux.HandleFunc("PUT /api/users", apiCfg.UpdateUserHandler)
	mux.HandleFunc("POST /api/login", apiCfg.LoginUserHandler)

	mux.HandleFunc("GET /api/chirps", apiCfg.GetAllChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.GetChirpHandler)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.DeleteChirpHandler)

	mux.HandleFunc("POST /api/refresh", apiCfg.RefreshTokenHandler)
	mux.HandleFunc("POST /api/revoke", apiCfg.RevokeRefreshTokenHandler)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.UpgradeUserHandler)

	// Example to serve to a url different from the directory folder names
	// dir := http.Dir("./assets/")
	// mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(dir)))

	log.Printf("Serving on port: %s\n", port)
	httpServer.ListenAndServe()
}
