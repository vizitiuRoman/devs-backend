package routes

import (
	"github.com/gorilla/mux"

	. "github.com/devs-backend/user-service/pkg/controllers"
	. "github.com/devs-backend/user-service/pkg/middlewares"
)

func InitRoutes() *mux.Router {
	router := mux.NewRouter()

	// Home
	router.HandleFunc("/", GetHome).Methods("GET")

	// Auth
	router.HandleFunc("/login", MiddlewareJSON(Login)).Methods("POST")
	router.HandleFunc("/register", MiddlewareJSON(Register)).Methods("POST")
	router.HandleFunc("/logout", MiddlewareJSON(Logout)).Methods("POST")

	// User
	router.HandleFunc("/users", MiddlewareAUTH(GetUsers)).Methods("GET")
	router.HandleFunc("/user/{id}", MiddlewareAUTH(GetUserByID)).Methods("GET")
	router.HandleFunc("/user/{id}", MiddlewareAUTH(UpdateUser)).Methods("PUT")
	router.HandleFunc("/user/{id}", MiddlewareAUTH(DeleteUserByID)).Methods("DELETE")

	return router
}
