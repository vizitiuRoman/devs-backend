package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// Private Func

func extractToken(r *http.Request) string {
	keys := r.URL.Query()
	token := keys.Get("token")
	if token != "" {
		return token
	}

	bearerToken := r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func pretty(data interface{}) {
	_, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println(err)
		return
	}
}

func prepareToken(extractedToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(extractedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return &jwt.Token{}, err
	}
	return token, nil
}

// Public Func

func CreateToken(userID uint32) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["userID"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 12).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("API_SECRET")))
}

func TokenValid(r *http.Request) error {
	extractedToken := extractToken(r)
	token, err := prepareToken(extractedToken)
	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		pretty(claims)
	}
	return nil
}

func EncodeToken(r *http.Request) (uint32, error) {
	extractedToken := extractToken(r)
	token, err := prepareToken(extractedToken)
	if err != nil {
		return 0, nil
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["userID"]), 10, 32)
		if err != nil {
			return 0, err
		}
		return uint32(uid), nil
	}
	return 0, nil
}
