package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	. "github.com/devsmd/pkg/auth"
	. "github.com/devsmd/pkg/db/models"
	. "github.com/devsmd/pkg/utils"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

const (
	userDataSession   = "user-session"
	userLoginAction   = "login"
	userDefaultAction = ""
)

var (
	key   = []byte(os.Getenv("SESSION_KEY"))
	store = sessions.NewCookieStore(key)
)

type response struct {
	UserID   uint32 `json:"userId"`
	Name     string `json:"name"`
	LastName string `json:"lastName"`
	Token    string `json:"token"`
}

// Private func

func prepareUser(r *http.Request, user *User) (*User, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return user, err
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		return user, err
	}
	return user, err
}

func saveToken(token string, userID uint32) error {
	tk := Token{Token: token, UserID: userID}
	tk.Prepare()
	_, err := tk.CreateOrUpdate(userID)
	if err != nil {
		return err
	}
	return nil
}

func setUserSession(session *sessions.Session, w http.ResponseWriter, r *http.Request, authenticated bool, userID uint32) {
	session.Values["authenticated"] = authenticated
	session.Values["userID"] = userID
	session.Save(r, w)
}

// Public func

func Login(w http.ResponseWriter, r *http.Request) {
	user, err := prepareUser(r, &User{})
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, errors.New(http.StatusText(http.StatusUnprocessableEntity)))
		return
	}

	user.Prepare()
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

	err = saveToken(token, receivedUser.ID)
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, errors.New(http.StatusText(http.StatusUnprocessableEntity)))
		return
	}

	session, _ := store.Get(r, userDataSession)
	setUserSession(session, w, r, true, receivedUser.ID)

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

	user.Prepare()
	err = user.Validate(userDefaultAction)
	if err != nil {
		ERROR(w, http.StatusBadRequest, errors.New(http.StatusText(http.StatusBadRequest)))
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

	err = saveToken(token, createdUser.ID)
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, errors.New(http.StatusText(http.StatusUnprocessableEntity)))
		return
	}

	session, _ := store.Get(r, userDataSession)
	setUserSession(session, w, r, true, createdUser.ID)

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, createdUser.ID))
	JSON(w, http.StatusCreated, response{
		UserID:   createdUser.ID,
		Name:     createdUser.Name,
		LastName: createdUser.LastName,
		Token:    token,
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, userDataSession)
	userID := session.Values["userID"]

	encoded, err := EncodeToken(r)
	if err != nil || encoded != userID {
		ERROR(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
		return
	}

	var token Token
	_, err = token.DeleteById(userID.(uint32))
	if err != nil {
		ERROR(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
		return
	}

	setUserSession(session, w, r, false, 0)

	JSON(w, http.StatusOK, true)
}
