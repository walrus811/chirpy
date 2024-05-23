package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/walrus811/chirpy/internal/database"
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
	const dbPath = "database.json"

	db, dbErr := database.NewDB(dbPath)
	if dbErr != nil {
		fmt.Println("Error creating database")
		return
	}

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
	mux.HandleFunc("GET /api/chirps/{chirpID}", func(w http.ResponseWriter, r *http.Request) {
		chirpID, err := strconv.Atoi(r.PathValue("chirpID"))

		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
			return
		}

		chirp, getErr := db.GetChirp(chirpID)
		if getErr != nil {
			respondWithError(w, http.StatusNotFound, "not found")
			return
		}

		respondWithJson(w, http.StatusOK, chirp)
	})
	mux.HandleFunc("GET /api/chirps", func(w http.ResponseWriter, r *http.Request) {
		chrips, err := db.GetChirps()

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Something went wrong")
			return
		}

		respondWithJson(w, http.StatusOK, chrips)
	})
	mux.HandleFunc("POST /api/chirps", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		reqObj := createChirpRequest{}
		err := decoder.Decode(&reqObj)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Something went wrong")
			return
		}

		if len(reqObj.Body) > 140 {
			respondWithError(w, http.StatusBadRequest, "Chirp is too long")
			return
		}

		newChirp, createErr := db.CreateChirp(reqObj.Body)

		if createErr != nil {
			respondWithError(w, http.StatusInternalServerError, "Something went wrong")
			return
		}

		respondWithJson(w, http.StatusCreated, newChirp)
	})

	mux.HandleFunc("POST /api/users", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		reqObj := createUserRequest{}
		err := decoder.Decode(&reqObj)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Something went wrong")
			return
		}

		newUser, createErr := db.CreateUser(reqObj.Email, reqObj.Password)

		if createErr != nil {
			respondWithError(w, http.StatusInternalServerError, "Something went wrong")
			return
		}

		resObj := createUserResponse{newUser.Id, newUser.Email}

		respondWithJson(w, http.StatusCreated, resObj)
	})

	mux.HandleFunc("POST /api/login", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		reqObj := loginUserRequest{}
		err := decoder.Decode(&reqObj)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Something went wrong")
			return
		}

		user, loginErr := db.LoginUser(reqObj.Email, reqObj.Password)

		if loginErr != nil {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		resObj := loginUserResponse{user.Id, user.Email}

		respondWithJson(w, http.StatusOK, resObj)
	})

	http.ListenAndServe(":"+port, mux)
}

type createChirpRequest struct {
	Body string `json:"body"`
}

type createUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type createUserResponse struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type loginUserResponse struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json;charset=UTF-8")
	json.NewEncoder(w).Encode(ErrorResponse{Error: msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json;charset=UTF-8")
	json.NewEncoder(w).Encode(payload)
}

func makeWordClean(word string) string {
	lower := strings.ToLower(word)
	if lower == "kerfuffle" || lower == "sharbert" || lower == "fornax" {
		return "****"
	}
	return word
}
