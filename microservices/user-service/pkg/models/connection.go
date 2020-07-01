package models

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	DB     *gorm.DB
	Client *redis.Client
)

func ConnectDB() {
	connectREDIS()
	connectPG(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"),
	)
	debug, _ := strconv.ParseBool(os.Getenv("DB_DEBUG"))
	DB.LogMode(debug)
	DB.AutoMigrate(&User{})
	fmt.Println("Debug mode:", debug)
}

func connectPG(DBDriver, DBUser, DBPassword, DBPort, DBHost, DBName string) {
	database, err := gorm.Open(DBDriver, fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		DBHost, DBPort, DBUser, DBName, DBPassword,
	))
	if err != nil {
		fmt.Println("Postgres cannot connect to", DBName)
		log.Fatal("Error", err)
	}
	fmt.Println("Postgres connect to", DBName)
	DB = database
}

func connectREDIS() {
	host, port := os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")
	if len(host) == 0 && len(port) == 0 {
		host, port = "127.0.0.1", "6379"
	}
	Client = redis.NewClient(&redis.Options{
		Addr: host + ":" + port,
	})
	_, err := Client.Ping().Result()
	if err != nil {
		fmt.Println("Redis cannot connect to", host+":"+port)
		log.Fatal("Error", err)
	}
	fmt.Println("Redis connect to", host+":"+port)
}
