package mvc

import (
	"fmt"
	"log"
	"net/http"
	"os"

	. "github.com/devs-backend/user-service/mvc/models"
	. "github.com/devs-backend/user-service/mvc/routes"
	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"
)

func Serve() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error load env", err)
	}
	fmt.Println("Load env")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	ConnectDB()
	routes, headers, methods, origins := InitRoutes()

	fmt.Println("App started", port)
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(headers, methods, origins)(routes)))
}