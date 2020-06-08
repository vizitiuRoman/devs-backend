package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	. "github.com/devsmd/pkg/auth"
	. "github.com/devsmd/pkg/models"
	. "github.com/devsmd/pkg/utils"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

const (
	userDataSession   = "user-data-session"
	userLoginAction   = "login"
	userDefaultAction = ""
)

var (
	key   = []byte(os.Getenv("SESSION_KEY"))
	store = sessions.NewCookieStore(key)
)

type response struct {
	UserId   uint32 `json:"userId"`
	Name     string `json:"name"`
	LastName string `json:"lastName"`
	Token    string `json:"token"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user.Prepare()
	err = user.Validate(userLoginAction)
	if err != nil {
		ERROR(w, http.StatusBadRequest, err)
		return
	}

	password := user.Password
	receivedUser, err := user.FindByEmail()
	if err != nil {
		ERROR(w, http.StatusNotFound, errors.New("User Not Found"))
		return
	}

	err = VerifyPassword(receivedUser.Password, password)
	if err != nil {
		ERROR(w, http.StatusBadRequest, errors.New("Wrong data"))
		return
	}

	token, err := CreateToken(receivedUser.ID)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	res := response{
		UserId:   receivedUser.ID,
		Name:     receivedUser.Name,
		LastName: receivedUser.LastName,
		Token:    token,
	}

	session, _ := store.Get(r, userDataSession)
	session.Values["authenticated"] = true
	session.Values["userId"] = receivedUser.ID
	session.Save(r, w)

	JSON(w, http.StatusOK, res)
}

func Register(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user.Prepare()
	err = user.Validate(userDefaultAction)
	if err != nil {
		ERROR(w, http.StatusBadRequest, err)
		return
	}

	createdUser, err := user.Create()
	if err != nil {
		ERROR(w, http.StatusInternalServerError, err)
		return
	}

	session, _ := store.Get(r, userDataSession)
	session.Values["authenticated"] = true
	session.Values["userId"] = createdUser.ID
	session.Save(r, w)

	token, err := CreateToken(createdUser.ID)
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	res := response{
		UserId:   createdUser.ID,
		Name:     createdUser.Name,
		LastName: createdUser.LastName,
		Token:    token,
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, createdUser.ID))
	JSON(w, http.StatusCreated, res)
}

func Logout(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, userDataSession)
	session.Values["authenticated"] = false
	session.Values["userId"] = nil
	session.Save(r, w)

	JSON(w, http.StatusOK, nil)
}
