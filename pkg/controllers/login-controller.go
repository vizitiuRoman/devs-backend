package controllers

import (
	"net/http"
	"os"

	. "github.com/devsmd/pkg/utils"
	"github.com/gorilla/sessions"
)

var (
	key   = []byte(os.Getenv("SESSION_KEY"))
	store = sessions.NewCookieStore(key)
)

func Login(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "auth-session")
	session.Values["authenticated"] = true
	session.Values["userId"] = 10
	session.Save(r, w)

	JSON(w, http.StatusOK, nil)
}

func Logout(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "auth-session")
	session.Values["authenticated"] = false
	session.Values["userId"] = nil
	session.Save(r, w)

	JSON(w, http.StatusOK, nil)
}

func Register(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "auth-session")
	session.Values["authenticated"] = true
	session.Values["userId"] = 10
	session.Save(r, w)

	JSON(w, http.StatusOK, nil)
}
