package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jental/freetesl-server/common"
	"github.com/jental/freetesl-server/db"
	"github.com/jental/freetesl-server/dtos"
)

var jwtKey = []byte("my_secret_key")
var tokens []string

type Claims struct {
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

	var responseDTO dtos.LoginResponseDTO
	valid := db.VerifyUser(request.Login, request.PasswordSha512)
	if valid {
		jwt, err := generateJWT(request.Login)
		if err != nil {
			log.Println(err)
			return
		}
		responseDTO = dtos.LoginResponseDTO{Valid: true, Token: &jwt}
	} else {
		responseDTO = dtos.LoginResponseDTO{Valid: false, Token: nil}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseDTO)
}

func generateJWT(login string) (string, error) {
	expirationTime := time.Now().Add(common.AUTH_JWT_EXPIRATION_TIME * time.Second)
	claims := &Claims{
		Username: login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtKey)
}
