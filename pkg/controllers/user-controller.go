package controllers

import (
	"errors"
	"fmt"
	"net/http"

	. "github.com/devsmd/pkg/auth"
	. "github.com/devsmd/pkg/db/models"
	. "github.com/devsmd/pkg/utils"
)

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, userDataSession)
	userID := session.Values["userID"]

	encoded, err := EncodeToken(r)
	if err != nil || encoded != userID {
		ERROR(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
		return
	}

	var user User
	err = user.DeleteById(userID.(uint32))
	if err != nil {
		ERROR(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
		return
	}

	w.Header().Set("Entity", fmt.Sprintf("%d", userID.(uint32)))
	JSON(w, http.StatusNoContent, "")
}
