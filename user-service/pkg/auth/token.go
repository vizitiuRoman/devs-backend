package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	. "github.com/devs-backend/user-service/pkg/models"
	. "github.com/devs-backend/user-service/pkg/utils"
	jwt "github.com/dgrijalva/jwt-go"
	. "github.com/twinj/uuid"
)

type AccessDetails struct {
	AccessUUID  string
	RefreshUUID string
	UserID      uint64
}

// Private func

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
	claims[Authorized] = true
	claims[UserID] = userID
	claims[AccessUUID] = accessUUID
	claims[RefreshUUID] = refreshUUID
	claims["exp"] = time.Now().Add(TokenExpires).Unix()
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv("API_SECRET")))
	if err != nil {
		return "", err
	}

	tokenDetails := &TokenDetails{
		AccessToken: token,
		AccessUUID:  accessUUID,
		RefreshUUID: refreshUUID,
		AtExpires:   time.Now().Add(AtExpires).Unix(),
		RtExpires:   time.Now().Add(RtExpires).Unix(),
	}
	err = tokenDetails.Create(userID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	extractedToken := extractToken(r)
	if extractedToken == "" {
		return &AccessDetails{}, errors.New("Can't extract token")
	}

	token, err := prepareToken(extractedToken)
	if err != nil {
		return &AccessDetails{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		accessUUID, ok := claims[AccessUUID].(string)
		if !ok {
			return &AccessDetails{}, err
		}
		refreshUUID, ok := claims[RefreshUUID].(string)
		if !ok {
			return &AccessDetails{}, err
		}
		userID, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["userID"]), 10, 32)
		if err != nil {
			return &AccessDetails{}, err
		}
		return &AccessDetails{
			AccessUUID:  accessUUID,
			RefreshUUID: refreshUUID,
			UserID:      userID,
		}, nil
	}
	return &AccessDetails{}, errors.New("Extract token error")
}

func FetchToken(accessDT *AccessDetails) (uint64, error) {
	token := TokenDetails{AccessUUID: accessDT.AccessUUID}
	userID, err := token.GetByUUID(accessDT.UserID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
