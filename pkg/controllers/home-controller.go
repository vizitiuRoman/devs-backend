package controllers

import (
	"net/http"

	. "github.com/devsmd/pkg/utils"
)

type User struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

func GetHome(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "auth-session")

	userId := session.Values["userId"]

	JSON(w, http.StatusOK, userId)
}
