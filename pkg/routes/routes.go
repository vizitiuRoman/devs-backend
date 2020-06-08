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

	router.HandleFunc("/home", MiddlewareAuth(GetHome)).Methods("GET")
	router.HandleFunc("/token", MiddlewareJSON(CheckToken)).Methods("GET")
	router.HandleFunc("/login", MiddlewareJSON(Login)).Methods("GET")

	return router, headers, methods, origins
}
