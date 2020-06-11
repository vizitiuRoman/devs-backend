package controllers

import (
	"net/http"

	. "github.com/devsmd/pkg/utils"
)

func GetHome(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, userDataSession)
	userID := session.Values["userID"]
	JSON(w, http.StatusOK, userID)
}
