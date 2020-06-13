package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	. "github.com/devs-backend/user-service/pkg/models"
	jwt "github.com/dgrijalva/jwt-go"
	. "github.com/twinj/uuid"
)

type AccessDetails struct {
	AccessUUID string
	UserID     uint64
}

// Private func

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

// Public func

func CreateToken(userID uint32) (string, error) {
	accessUUID := NewV4().String()
	refreshUUID := NewV4().String()
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["userID"] = userID
	claims["accessUUID"] = accessUUID
	claims["refreshUUID"] = refreshUUID
	claims["exp"] = time.Now().Add(time.Hour * 12).Unix()
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv("API_SECRET")))
	if err != nil {
		return "", err
	}

	tokenDetails := &TokenDetails{
		AccessToken: token,
		AccessUUID:  accessUUID,
		RefreshUUID: refreshUUID,
		AtExpires:   time.Now().Add(time.Hour * 12).Unix(),
		RtExpires:   time.Now().Add(time.Hour * 24 * 7).Unix(),
	}
	err = tokenDetails.Create(userID)
	if err != nil {
		return "", err
	}
	return token, nil
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

func EncodeToken(r *http.Request) (*AccessDetails, error) {
	extractedToken := extractToken(r)
	if extractedToken == "" {
		return &AccessDetails{}, errors.New("Extract token")
	}

	token, err := prepareToken(extractedToken)
	if err != nil {
		return &AccessDetails{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		accessUuid, ok := claims["accessUUID"].(string)
		if !ok {
			return &AccessDetails{}, err
		}
		userID, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["userID"]), 10, 32)
		if err != nil {
			return &AccessDetails{}, err
		}
		return &AccessDetails{
			AccessUUID: accessUuid,
			UserID:     userID,
		}, nil
	}
	return &AccessDetails{}, errors.New("Extract token")
}

func GetToken(accessD *AccessDetails) (uint64, error) {
	userid, err := Client.Get(accessD.AccessUUID).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseUint(userid, 10, 64)
	if accessD.UserID != userID {
		return 0, errors.New("unauthorized")
	}
	return userID, nil
}
