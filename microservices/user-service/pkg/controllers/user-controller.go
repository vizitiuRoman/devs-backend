package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	. "github.com/devs-backend/user-service/pkg/auth"
	. "github.com/devs-backend/user-service/pkg/models"
	. "github.com/devs-backend/user-service/pkg/utils"
	"github.com/gorilla/mux"
)

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		fmt.Println(err)
		ERROR(w, http.StatusUnprocessableEntity, errors.New(http.StatusText(http.StatusUnprocessableEntity)))
		return
	}

	user := User{ID: uint32(userID)}
	receivedUser, err := user.FindByID()
	if err != nil {
		ERROR(w, http.StatusNotFound, errors.New(http.StatusText(http.StatusNotFound)))
		return
	}

	JSON(w, http.StatusOK, receivedUser)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	var user User
	receivedUsers, err := user.FindAll()
	if err != nil {
		ERROR(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
		return
	}
	JSON(w, http.StatusOK, receivedUsers)
}

func DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, errors.New(http.StatusText(http.StatusUnprocessableEntity)))
		return
	}

	extractedToken, err := ExtractTokenMetadata(r)
	if err != nil {
		ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	if userID != extractedToken.UserID {
		ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	user := User{ID: uint32(userID)}
	err = user.DeleteByID(extractedToken.AccessUUID, extractedToken.RefreshUUID)
	if err != nil {
		ERROR(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
		return
	}

	JSON(w, http.StatusOK, true)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, errors.New(http.StatusText(http.StatusUnprocessableEntity)))
		return
	}

	user, err := prepareUser(r, &User{})
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, errors.New(http.StatusText(http.StatusUnprocessableEntity)))
		return
	}

	err = user.Validate(UPDATE)
	if err != nil {
		ERROR(w, http.StatusBadRequest, err)
		return
	}

	extractedToken, err := ExtractTokenMetadata(r)
	if err != nil {
		ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	user.ID = uint32(userID)
	if userID != extractedToken.UserID {
		ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	updatedUser, err := user.Update()
	if err != nil {
		ERROR(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
		return
	}

	JSON(w, http.StatusOK, updatedUser)
}
