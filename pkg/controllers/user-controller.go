package controllers

import (
	"fmt"
	"net/http"

	. "github.com/devsmd/pkg/auth"
	. "github.com/devsmd/pkg/models"
	. "github.com/devsmd/pkg/utils"
)

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, userDataSession)
	userID := session.Values["userID"]

	var user User
	_, err := user.DeleteById(userID.(uint32))
	if err != nil {
		ERROR(w, http.StatusInternalServerError, err)
		return
	}

	encoded, _ := EncodeToken(r)
	fmt.Println(encoded)

	w.Header().Set("Entity", fmt.Sprintf("%d", userID.(uint32)))
	JSON(w, http.StatusNoContent, "")
}
