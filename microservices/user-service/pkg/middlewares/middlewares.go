package middlewares

import (
	"errors"
	"log"
	"net/http"

	. "github.com/devs-backend/user-service/pkg/auth"
	. "github.com/devs-backend/user-service/pkg/utils"
)

func MiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func MiddlewareAUTH(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		token, err := ExtractTokenMetadata(r)
		if err != nil {
			log.Println("MiddlewareAUTH-ExtractTokenMetadata:", err)
			ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
			return
		}

		_, err = FetchToken(&AccessDetails{
			AccessUUID: token.AccessUUID,
			UserID:     token.UserID,
		})
		if err != nil {
			log.Println("MiddlewareAUTH-FetchToken:", err)
			ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
			return
		}
		next(w, r)
	}
}
