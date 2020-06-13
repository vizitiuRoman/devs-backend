package middlewares

import (
	"errors"
	"fmt"
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
		token, err := EncodeToken(r)
		if err != nil {
			ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
			return
		}

		_, err = GetToken(&AccessDetails{
			AccessUUID: token.AccessUUID,
			UserID:     token.UserID,
		})
		if err != nil {
			ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
			return
		}

		fmt.Println("Q")

		next(w, r)
	}
}
