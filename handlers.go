package main

import (
	"fmt"
	"net/http"
	"saurabh/chirpy.com/m/internal/database"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	db             *database.Queries
	platform       string
}

func (cfg *apiConfig) MiddlewareMetricsIncrement(next http.Handler) http.Handler {

	InnerFunc := func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(InnerFunc)
}

func (cfg *apiConfig) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("ContentType", "text/html; charset=utf-8")
	w.WriteHeader(200)
	htmlContent := fmt.Sprintf(`<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
	</html>`, cfg.fileServerHits.Load())

	w.Write([]byte(htmlContent))

}

func (cfg *apiConfig) ResetMetricsHandler(w http.ResponseWriter, r *http.Request) {

	cfg.fileServerHits.Store(0)

	if cfg.platform == "dev" {
		cfg.db.DeleteAllUsers(r.Context())
		cfg.db.DeleteAllChirps(r.Context())
	}

	w.Header().Add("ContentType", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}
