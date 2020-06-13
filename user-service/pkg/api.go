package pkg

import (
	"fmt"
	"log"
	"net/http"
	"os"

	. "github.com/devs-backend/user-service/pkg/routes"
)

func Serve() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	routes := InitRoutes()

	fmt.Println("App started", port)
	log.Fatal(http.ListenAndServe(":"+port, routes))
}
