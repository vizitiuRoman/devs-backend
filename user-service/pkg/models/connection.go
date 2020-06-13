package models

import (
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

var (
	db     *gorm.DB
	Client *redis.Client
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error load env", err)
	} else {
		fmt.Println("Load env")
	}

	connectREDIS()
	connectPG(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"),
	)

	db.AutoMigrate(&User{})
}

func connectPG(DbDriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		DbHost, DbPort, DbUser, DbName, DbPassword,
	)

	database, err := gorm.Open(DbDriver, DBURL)
	if err != nil {
		fmt.Println("Postgres can't connect to", DbName)
		log.Fatal("Error", err)
	} else {
		fmt.Println("Connect to", DbName)
	}
	db = database
}

func connectREDIS() {
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	Client = redis.NewClient(&redis.Options{
		Addr: dsn,
	})
	_, err := Client.Ping().Result()
	if err != nil {
		fmt.Println("Redis can't connect to", dsn)
		log.Fatal("Error", err)
	}
}
