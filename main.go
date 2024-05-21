package main

import (
	"fmt"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	cfg := &apiConfig{fileserverHits: 0}
	mux := http.NewServeMux()

	mux.Handle("GET /app/*", http.StripPrefix("/app/", cfg.middlewareMetricsInc(http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "text/plain;charset=UTF-8")
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("GET /api/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "text/plain;charset=UTF-8")
		w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits)))
	})
	mux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "text/html")
		w.Write([]byte(fmt.Sprintf(`<html>
		<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
		</body>
		</html>`, cfg.fileserverHits)))
	})
	mux.HandleFunc("GET /api/reset", func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits = 0
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "text/plain;charset=UTF-8")
		w.Write([]byte("Hits reset to 0"))
	})

	http.ListenAndServe(":8080", mux)
}
