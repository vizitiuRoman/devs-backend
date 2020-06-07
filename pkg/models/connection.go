package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var db *gorm.DB

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error load env", err)
	} else {
		fmt.Println("Load env")
	}

	connect(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"),
	)
	db.AutoMigrate()
}

func connect(DbDriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		DbHost, DbPort, DbUser, DbName, DbPassword,
	)

	database, err := gorm.Open(DbDriver, DBURL)
	if err != nil {
		fmt.Println("Can't connect to", DbName)
		log.Fatal("Error", err)
	} else {
		fmt.Println("Connect to", DbName)
	}
	db = database
}
