package middlewares

import (
	"errors"
	"net/http"

	. "github.com/devsmd/pkg/auth"
	. "github.com/devsmd/pkg/utils"
)

func MiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func MiddlewareAUTH(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := TokenValid(r)
		if err != nil {
			ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}
