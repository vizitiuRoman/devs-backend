package routes

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	. "github.com/devs-backend/user-service/pkg/controllers"
	. "github.com/devs-backend/user-service/pkg/middlewares"
)

func InitRoutes() (*mux.Router, handlers.CORSOption, handlers.CORSOption, handlers.CORSOption) {
	router := mux.NewRouter()

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "token"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})

	// Home
	router.HandleFunc("/", MiddlewareAUTH(GetHome)).Methods("GET")

	// Auth
	router.HandleFunc("/login", MiddlewareJSON(Login)).Methods("POST")
	router.HandleFunc("/register", MiddlewareJSON(Register)).Methods("POST")
	router.HandleFunc("/logout", MiddlewareJSON(Logout)).Methods("POST")

	// User
	router.HandleFunc("/user", MiddlewareAUTH(GetUsers)).Methods("GET")
	router.HandleFunc("/user/{id}", MiddlewareAUTH(GetUserByID)).Methods("GET")
	router.HandleFunc("/user", MiddlewareAUTH(UpdateUser)).Methods("POST")
	router.HandleFunc("/user/{id}", MiddlewareAUTH(DeleteUserByID)).Methods("DELETE")

	return router, headers, methods, origins
}
