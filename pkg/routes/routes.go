package routes

import (
	. "github.com/devsmd/pkg/controllers"
	. "github.com/devsmd/pkg/middlewares"
	. "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func InitRoutes() (*mux.Router, CORSOption, CORSOption, CORSOption) {
	router := mux.NewRouter()

	headers := AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	origins := AllowedOrigins([]string{"*"})

	// Home Router
	router.HandleFunc("/api/home", MiddlewareAUTH(GetHome)).Methods("GET")

	// Auth Router
	router.HandleFunc("/login", MiddlewareJSON(Login)).Methods("POST")
	router.HandleFunc("/register", MiddlewareJSON(Register)).Methods("POST")
	router.HandleFunc("/logout", MiddlewareJSON(Logout)).Methods("POST")

	// User Router
	router.HandleFunc("/api/user/delete", MiddlewareAUTH(DeleteUser)).Methods("DELETE")

	return router, headers, methods, origins
}
