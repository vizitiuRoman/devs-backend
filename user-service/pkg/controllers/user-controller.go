package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	. "github.com/devs-backend/user-service/pkg/auth"
	. "github.com/devs-backend/user-service/pkg/models"
	. "github.com/devs-backend/user-service/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

const (
	userLoginAction   = "login"
	userDefaultAction = ""
)

var key = []byte(os.Getenv("SESSION_KEY"))

type response struct {
	UserID   uint32 `json:"userId"`
	Name     string `json:"name"`
	LastName string `json:"lastName"`
	Token    string `json:"token"`
}

func prepareUser(r *http.Request, user *User) (*User, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return user, err
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		return user, err
	}
	user.Prepare()
	return user, err
}

func Login(w http.ResponseWriter, r *http.Request) {
	user, err := prepareUser(r, &User{})
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, errors.New(http.StatusText(http.StatusUnprocessableEntity)))
		return
	}

	err = user.Validate(userLoginAction)
	if err != nil {
		ERROR(w, http.StatusBadRequest, errors.New(http.StatusText(http.StatusBadRequest)))
		return
	}

	password := user.Password
	receivedUser, err := user.FindByEmail()
	if err != nil {
		ERROR(w, http.StatusNotFound, errors.New(http.StatusText(http.StatusNotFound)))
		return
	}

	err = VerifyPassword(receivedUser.Password, password)
	if err != nil {
		ERROR(w, http.StatusBadRequest, errors.New(http.StatusText(http.StatusBadRequest)))
		return
	}

	token, err := CreateToken(receivedUser.ID)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		ERROR(w, http.StatusUnprocessableEntity, errors.New(http.StatusText(http.StatusUnprocessableEntity)))
		return
	}

	JSON(w, http.StatusOK, response{
		UserID:   receivedUser.ID,
		Name:     receivedUser.Name,
		LastName: receivedUser.LastName,
		Token:    token,
	})
}

func Register(w http.ResponseWriter, r *http.Request) {
	user, err := prepareUser(r, &User{})
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, errors.New(http.StatusText(http.StatusUnprocessableEntity)))
		return
	}

	err = user.Validate(userDefaultAction)
	if err != nil {
		ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	createdUser, err := user.Create()
	if err != nil {
		ERROR(w, http.StatusConflict, errors.New(http.StatusText(http.StatusConflict)))
		return
	}

	token, err := CreateToken(createdUser.ID)
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, errors.New(http.StatusText(http.StatusUnprocessableEntity)))
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, createdUser.ID))
	JSON(w, http.StatusCreated, response{
		UserID:   createdUser.ID,
		Name:     createdUser.Name,
		LastName: createdUser.LastName,
		Token:    token,
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	_, err := EncodeToken(r)
	if err != nil {
		ERROR(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
		return
	}
	JSON(w, http.StatusOK, true)
}
