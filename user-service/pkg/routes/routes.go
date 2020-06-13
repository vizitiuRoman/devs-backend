package routes

import (
	"github.com/gorilla/mux"

	. "github.com/devs-backend/user-service/pkg/controllers"
	. "github.com/devs-backend/user-service/pkg/middlewares"
)

func InitRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/login", MiddlewareJSON(Login)).Methods("POST")
	router.HandleFunc("/register", MiddlewareJSON(Register)).Methods("POST")
	router.HandleFunc("/logout", MiddlewareJSON(Logout)).Methods("POST")

	return router
}
