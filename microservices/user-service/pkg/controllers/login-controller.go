package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	. "github.com/devs-backend/user-service/pkg/auth"
	. "github.com/devs-backend/user-service/pkg/models"
	. "github.com/devs-backend/user-service/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

const (
	userLoginAction   = "login"
	userDefaultAction = ""
)

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
		ERROR(w, http.StatusBadRequest, err)
		return
	}

	password := user.Password
	receivedUser, err := user.FindByEmail()
	if err != nil {
		ERROR(w, http.StatusBadRequest, errors.New("Invalid Email Or Password"))
		return
	}

	err = VerifyPassword(receivedUser.Password, password)
	if err != nil {
		ERROR(w, http.StatusBadRequest, errors.New("Invalid Email Or Password"))
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
		ERROR(w, http.StatusBadRequest, err)
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
	extractedToken, err := ExtractTokenMetadata(r)
	if err != nil {
		ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	token := TokenDetails{
		AccessUUID:  extractedToken.AccessUUID,
		RefreshUUID: extractedToken.RefreshUUID,
	}
	err = token.DeleteByUUID()
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, errors.New(http.StatusText(http.StatusUnprocessableEntity)))
		return
	}

	JSON(w, http.StatusOK, true)
}
