package pkg

import (
	"fmt"
	. "github.com/devsmd/pkg/routes"
	. "github.com/gorilla/handlers"
	"log"
	"net/http"
	"os"
)

func Serve() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	routes, headers, methods, origins := InitRoutes()

	fmt.Println("App started", port)
	log.Fatal(http.ListenAndServe(":"+port, CORS(headers, methods, origins)(routes)))
}
