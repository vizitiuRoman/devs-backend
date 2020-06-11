package controllers

import (
	"errors"
	"net/http"

	. "github.com/devsmd/pkg/auth"
	. "github.com/devsmd/pkg/utils"
)

func GetHome(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, userDataSession)
	userID := session.Values["userID"]

	_, err := EncodeToken(r)
	if err != nil {
		ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	JSON(w, http.StatusOK, userID)
}
