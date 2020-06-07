package routes

import (
	. "github.com/devsmd/pkg/controllers"
	. "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func InitRoutes() (*mux.Router, CORSOption, CORSOption, CORSOption) {
	router := mux.NewRouter()

	headers := AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	origins := AllowedOrigins([]string{"*"})

	router.HandleFunc("/home", GetHome).Methods("GET")

	return router, headers, methods, origins
}
