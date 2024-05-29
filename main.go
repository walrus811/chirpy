package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/walrus811/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits int
	jwtSecret      string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}
	const filepathRoot = "."
	const port = "8080"
	const dbPath = "database.json"

	db, dbErr := database.NewDB(dbPath)
	if dbErr != nil {
		fmt.Println("Error creating database")
		return
	}

	cfg := &apiConfig{fileserverHits: 0, jwtSecret: os.Getenv("JWT_SECRET")}
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
		authKey := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]

		jwtClaim, getJwtClainErr := getJWTClaim(cfg.jwtSecret, authKey)

		if getJwtClainErr != nil {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		userIdStr, getSubjectErr := jwtClaim.GetSubject()

		if getSubjectErr != nil {
			respondWithError(w, http.StatusInternalServerError, "Something went wrong")
			return
		}

		userId, atoiErr := strconv.Atoi(userIdStr)

		if atoiErr != nil {
			respondWithError(w, http.StatusInternalServerError, "Something went wrong")
			return
		}

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

		newChirp, createErr := db.CreateChirp(reqObj.Body, userId)

		if createErr != nil {
			respondWithError(w, http.StatusInternalServerError, "Something went wrong")
			return
		}

		respondWithJson(w, http.StatusCreated, newChirp)
	})

	
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", func(w http.ResponseWriter, r *http.Request) {
		authKey := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]

		jwtClaim, getJwtClainErr := getJWTClaim(cfg.jwtSecret, authKey)

		if getJwtClainErr != nil {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		userIdStr, getSubjectErr := jwtClaim.GetSubject()

		if getSubjectErr != nil {
			respondWithError(w, http.StatusInternalServerError, "Something went wrong")
			return
		}

		userId, atoiErr := strconv.Atoi(userIdStr)

		if atoiErr != nil {
			respondWithError(w, http.StatusInternalServerError, "Something went wrong")
			return
		}

		chirpID, getChirpIDErr := strconv.Atoi(r.PathValue("chirpID"))

		if getChirpIDErr != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
			return
		}


		_, getChirpErr := db.GetChirp(chirpID)
	
		if getChirpErr != nil {
			respondWithError(w, http.StatusForbidden, "Forbidden")
			return
		}

		deleteErr := db.DeleteChirp(chirpID)

		if deleteErr != nil {
			respondWithError(w, http.StatusInternalServerError, "Something went wrong")
			return
		}

		respondWithJson(w, http.StatusNoContent, nil)
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

		token, tokenErr := getJWTString(cfg.jwtSecret, strconv.Itoa(user.Id))

		if tokenErr != nil {
			respondWithError(w, http.StatusInternalServerError, "Something went wrong")
			return
		}

		refreshToken, CreateRefreshTokenErr := db.CreateRefreshToken(user.Id)

		if CreateRefreshTokenErr != nil {
			respondWithError(w, http.StatusInternalServerError, "Something went wrong")
			return
		}

		resObj := loginUserResponse{user.Id, user.Email, token, refreshToken}

		respondWithJson(w, http.StatusOK, resObj)
	})

	mux.HandleFunc("PUT /api/users", func(w http.ResponseWriter, r *http.Request) {
		authKey := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]

		jwtClaim, getJwtClainErr := getJWTClaim(cfg.jwtSecret, authKey)

		if getJwtClainErr != nil {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		userIdStr, getSubjectErr := jwtClaim.GetSubject()

		if getSubjectErr != nil {
			fmt.Println(getSubjectErr)
			respondWithError(w, http.StatusInternalServerError, "Unauthorized")
			return
		}

		userId, atoiErr := strconv.Atoi(userIdStr)

		if atoiErr != nil {
			respondWithError(w, http.StatusInternalServerError, "Something went wrong")
			return
		}

		decoder := json.NewDecoder(r.Body)
		reqObj := updateUserRequest{}
		err := decoder.Decode(&reqObj)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Something went wrong")
			return
		}

		user, updateErr := db.UpdateUser(userId, reqObj.Email, reqObj.Password)

		if updateErr != nil {
			respondWithError(w, http.StatusInternalServerError, "Something went wrong")
			return
		}

		resObj := updateUserResponse{user.Id, user.Email}

		respondWithJson(w, http.StatusOK, resObj)
	})

	mux.HandleFunc("POST /api/refresh", func(w http.ResponseWriter, r *http.Request) {
		refershToken := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]

		userId, getUserIdErr := db.GetUserIdByToken(refershToken)

		if getUserIdErr != nil {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		newToken, getTokenErr := getJWTString(cfg.jwtSecret, strconv.Itoa(userId))

		if getTokenErr != nil {
			respondWithError(w, http.StatusInternalServerError, "Something went wrong")
			return
		}

		resObj := refershTokrnResponse{newToken}

		respondWithJson(w, http.StatusOK, resObj)
	})

	mux.HandleFunc("POST /api/revoke", func(w http.ResponseWriter, r *http.Request) {
		refershToken := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]

		deleteErr := db.DeleteRefreshToken(refershToken)

		if deleteErr != nil {
			respondWithError(w, http.StatusInternalServerError, "Something went wrong")
			return
		}

		w.WriteHeader(http.StatusNoContent)
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

type updateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type createUserResponse struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type updateUserResponse struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type loginUserResponse struct {
	Id           int    `json:"id"`
	Email        string `json:"email"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type refershTokrnResponse struct {
	Token string `json:"token"`
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

func getJWTString(signKey, id string) (string, error) {
	claim := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		Subject:   id,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(1) * time.Hour)),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(signKey))

	return token, err
}

func getJWTClaim(signKey, token string) (jwt.Claims, error) {
	var jwtClaim jwt.Claims = &jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(token, jwtClaim, func(token *jwt.Token) (interface{}, error) {
		return []byte(signKey), nil
	})

	return jwtClaim, err
}
