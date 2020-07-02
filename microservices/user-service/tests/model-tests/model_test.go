package model_tests

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"

	. "github.com/devs-backend/user-service/pkg/models"
	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	err := godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}

	connectREDIS()
	connectPG(os.Getenv("TEST_DB_DRIVER"), os.Getenv("TEST_DB_USER"), os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_PORT"), os.Getenv("TEST_DB_HOST"), os.Getenv("TEST_DB_NAME"),
	)
	debug, _ := strconv.ParseBool(os.Getenv("TEST_DB_DEBUG"))
	DB.LogMode(debug)
	os.Exit(m.Run())
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

func refreshUserTable() error {
	err := DB.DropTableIfExists(User{}).Error
	if err != nil {
		return err
	}
	err = DB.AutoMigrate(User{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed user table")
	return nil
}

func seedOneUser() (User, error) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user := User{
		Email:    "devs@gmail.com",
		Password: "password",
		Name:     "pet",
		LastName: "pets",
	}
	err = DB.Model(&User{}).Create(&user).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func seedUsers() ([]User, error) {
	users := []User{
		User{
			Email:    "gopher@gmail.com",
			Password: "password",
			Name:     "gopher",
			LastName: "chitaica",
		},
		User{
			Email:    "roma@gmail.com",
			Password: "password",
			Name:     "roma",
			LastName: "romanov",
		},
	}
	for i, _ := range users {
		err := DB.Model(&User{}).Create(&users[i]).Error
		if err != nil {
			return []User{}, err
		}
	}
	return users, nil
}
