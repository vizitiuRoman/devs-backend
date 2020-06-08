package controllers

import (
	"fmt"
	"net/http"

	. "github.com/devsmd/pkg/auth"
	. "github.com/devsmd/pkg/utils"
)

func GetHome(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, userDataSession)
	userID := session.Values["userID"]

	encode, _ := EncodeToken(r)
	fmt.Println(encode)

	JSON(w, http.StatusOK, userID)
}
