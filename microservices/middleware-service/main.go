package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var Client *redis.Client

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error load env", err)
	}
	fmt.Println("Load env")

	host, redisPort := os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")
	if len(host) == 0 && len(redisPort) == 0 {
		host, redisPort = "127.0.0.1", "6379"
	}
	Client = redis.NewClient(&redis.Options{
		Addr: host + ":" + redisPort,
	})
	_, err = Client.Ping().Result()
	if err != nil {
		fmt.Println("Redis can't connect to", host+":"+redisPort)
		log.Fatal("Error", err)
	}
	fmt.Println("Redis connect to", host+":"+redisPort)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8070"
	}

	router := mux.NewRouter()
	router.HandleFunc("/", handler).Methods("GET")

	fmt.Println("Middleware-service started", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	token, err := ExtractTokenMetadata(r)
	if err != nil {
		ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	_, err = FetchToken(&AccessDetails{
		AccessUUID: token.AccessUUID,
		UserID:     token.UserID,
	})
	if err != nil {
		ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	JSON(w, http.StatusOK, true)
}

func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}

func ERROR(w http.ResponseWriter, statusCode int, err error) {
	if err != nil {
		JSON(w, statusCode, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	JSON(w, http.StatusBadRequest, err)
}

const (
	UserID      = "userID"
	AccessUUID  = "accessUUID"
	RefreshUUID = "refreshUUID"
)

type AccessDetails struct {
	AccessToken string
	AccessUUID  string
	RefreshUUID string
	UserID      uint64
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
			return &AccessDetails{}, errors.New("Can't get accessUUID")
		}
		refreshUUID, ok := claims[RefreshUUID].(string)
		if !ok {
			return &AccessDetails{}, errors.New("Can't get refreshUUID")
		}
		userID, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims[UserID]), 10, 32)
		if err != nil {
			return &AccessDetails{}, errors.New("Can't get userID")
		}
		return &AccessDetails{
			AccessUUID:  accessUUID,
			RefreshUUID: refreshUUID,
			UserID:      userID,
		}, nil
	}
	return &AccessDetails{}, errors.New("ExtractTokenMetadata error")
}

func FetchToken(accessDT *AccessDetails) (uint64, error) {
	userid, err := Client.Get(accessDT.AccessUUID).Result()
	if err != nil {
		return 0, err
	}

	userID, _ := strconv.ParseUint(userid, 10, 64)
	if accessDT.UserID != userID {
		return 0, errors.New(http.StatusText(http.StatusUnauthorized))
	}
	return userID, nil
}
