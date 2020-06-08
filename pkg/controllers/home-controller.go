package controllers

import (
	"net/http"

	. "github.com/devsmd/pkg/utils"
)

func GetHome(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, userDataSession)

	userId := session.Values["userId"]

	JSON(w, http.StatusOK, userId)
}
