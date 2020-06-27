package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis/v7"
	"github.com/joho/godotenv"
)

var Client *redis.Client

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error load env", err)
	}
	fmt.Println("Load env")

	host, redisPort := os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")
	if len(host) == 0 && len(redisPort) == 0 {
		host, redisPort = "127.0.0.1", "6379"
	}
	Client = redis.NewClient(&redis.Options{
		Addr: host + ":" + redisPort,
	})
	_, err = Client.Ping().Result()
	if err != nil {
		fmt.Println("Redis can't connect to", host+":"+redisPort)
		log.Fatal("Error", err)
	}
	fmt.Println("Redis connect to", host+":"+redisPort)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8070"
	}

	routes := http.NewServeMux()
	routes.HandleFunc("/", handler)

	log.Fatal(http.ListenAndServe(":"+port, routes))
}
