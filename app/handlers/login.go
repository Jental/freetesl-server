package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jental/freetesl-server/common"
	"github.com/jental/freetesl-server/db"
	"github.com/jental/freetesl-server/dtos"
)

var jwtKey = []byte("jkh7hlkjg56'jkl")

type Claims struct {
	UserID   int    `json:"userid"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func Login(w http.ResponseWriter, req *http.Request) {
	var decoder = json.NewDecoder(req.Body)
	var request dtos.LoginDTO
	err := decoder.Decode(&request)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("Login: [%s]: password: '%s'\n", request.Login, request.PasswordSha512)

	var responseDTO dtos.LoginResponseDTO
	valid, userID := db.VerifyUser(request.Login, request.PasswordSha512)
	if valid {
		jwt, err := generateJWT(request.Login, *userID)
		if err != nil {
			log.Println(err)
			return
		}
		responseDTO = dtos.LoginResponseDTO{Valid: true, Token: &jwt}
	} else {
		log.Printf("Login: [%s]: failed", request.Login)
		responseDTO = dtos.LoginResponseDTO{Valid: false, Token: nil}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseDTO)
}

func generateJWT(login string, userID int) (string, error) {
	expirationTime := time.Now().Add(common.AUTH_JWT_EXPIRATION_TIME * time.Second)
	claims := &Claims{
		UserID:   userID,
		Username: login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtKey)
}

func verifyJWT(token string) (bool, *int, int) {
	claims := &Claims{}
	parsed, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return false, nil, http.StatusUnauthorized
		} else {
			return false, nil, http.StatusBadRequest
		}
	}

	if !parsed.Valid {
		return false, nil, http.StatusUnauthorized
	}

	return true, &claims.UserID, http.StatusOK
}

func AuthCheckMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Printf("AuthCheckMiddleware: req: %s", req.URL)

		fullAuthHeader := req.Header.Get("Authorization")
		var valid bool = false
		var statusCode int = http.StatusUnauthorized
		var userID *int = nil
		if fullAuthHeader != "" {
			parts := strings.Split(fullAuthHeader, " ")
			if len(parts) > 1 {
				token := parts[1]
				valid, userID, statusCode = verifyJWT(token)
			}
		}

		if !valid {
			log.Printf("AuthCheckMiddleware: req: %s: unathorized", req.URL)
			w.WriteHeader(statusCode)
			return
		}

		log.Printf("AuthCheckMiddleware: req: %s: userID: %d", req.URL, *userID)
		ctx := context.WithValue(req.Context(), "userID", *userID)
		newReq := req.WithContext(ctx)

		next.ServeHTTP(w, newReq)
	})
}
