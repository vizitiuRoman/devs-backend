package modeltests

import (
	"fmt"
	"log"
	"os"
	"testing"

	. "github.com/devs-backend/user-service/pkg/models"
	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	err := godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}

	connectREDIS()
	connectPG(os.Getenv("TEST_DB_DRIVER"), os.Getenv("TEST_DB_USER"), os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_PORT"), os.Getenv("TEST_DB_HOST"), os.Getenv("TEST_DB_NAME"),
	)

	os.Exit(m.Run())
}

func connectPG(TESTDBDriver, TESTDBUser, TESTDBPassword, TESTDBPort, TESTDBHost, TESTDBName string) {
	DBURL := fmt.Sprintf("host=%s port=%s user=%s name=%s sslmode=disable password=%s",
		TESTDBHost, TESTDBPort, TESTDBUser, TESTDBName, TESTDBPassword,
	)

	database, err := gorm.Open(TESTDBDriver, DBURL)
	if err != nil {
		fmt.Println("Postgres can't connect to", TESTDBName)
		log.Fatal("Error", err)
	}
	fmt.Println("Postgres connect to", TESTDBName)
	db = database
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
		fmt.Println("Redis can't connect to", host+":"+port)
		log.Fatal("Error", err)
	}
	fmt.Println("Redis connect to", host+":"+port)
}

func refreshUserTable() error {
	err := db.DropTableIfExists(User{}).Error
	if err != nil {
		return err
	}
	err = db.AutoMigrate(User{}).Error
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
		Email:    "pet@gmail.com",
		Password: "password",
		Name:     "pet",
		LastName: "pets",
	}
	err = db.Model(&User{}).Create(&user).Error
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
		err := db.Model(&User{}).Create(&users[i]).Error
		if err != nil {
			return []User{}, err
		}
	}
	return users, nil
}

func refreshUserAndPostTable() error {
	err := db.DropTableIfExists(&User{}).Error
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&User{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed tables")
	return nil
}
