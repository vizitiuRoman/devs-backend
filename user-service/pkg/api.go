package pkg

import (
	"fmt"
	"log"
	"net/http"
	"os"

	. "github.com/devs-backend/user-service/pkg/models"
	. "github.com/devs-backend/user-service/pkg/routes"
	"github.com/joho/godotenv"
)

func Serve() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error load env", err)
	}
	fmt.Println("Load env")

	ConnectDB()
	routes := InitRoutes()

	fmt.Println("App started", port)
	log.Fatal(http.ListenAndServe(":"+port, routes))
}
