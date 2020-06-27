package modeltests

import (
	"fmt"
	"log"
	"testing"

	. "github.com/devs-backend/user-service/pkg/models"
	"gopkg.in/stretchr/testify.v1/assert"
)

func TestCreateUser(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatalf("Error refreshUserTable: %v", err)
	}

	user := User{
		ID:       1,
		Name:     "Devs",
		LastName: "Devs",
		Email:    "devs@mail.ru",
		Password: "password",
	}
	createdUser, err := user.Create()
	if err != nil {
		t.Errorf("User Create: %v", err)
		return
	}

	assert.Equal(t, createdUser.ID, user.ID)
	assert.Equal(t, createdUser.Name, user.Name)
	assert.Equal(t, createdUser.LastName, user.LastName)
	assert.Equal(t, createdUser.Email, user.Email)
}

func TestUpdateUser(t *testing.T) {
	_, err := seedOneUser()
	if err != nil {
		fmt.Printf("Error seedOneUser: %v", err)
	}

	user := User{
		ID:       1,
		Name:     "Updated",
		LastName: "Updated",
	}
	updatedUser, err := user.Update()
	if err != nil {
		t.Errorf("User Update: %v", err)
		return
	}

	assert.Equal(t, updatedUser.ID, user.ID)
	assert.Equal(t, updatedUser.Name, user.Name)
	assert.Equal(t, updatedUser.LastName, user.LastName)
}

func TestDeleteUserByID(t *testing.T) {
	seededUser, err := seedOneUser()
	if err != nil {
		fmt.Printf("Error seedOneUser: %v", err)
	}

	tokenDetails := &TokenDetails{
		AccessUUID:  "123-22-444",
		RefreshUUID: "123-5343-333",
	}
	err = tokenDetails.Create(seededUser.ID)
	if err != nil {
		t.Errorf("Error Create tokenDetails: %v", err)
	}

	user := User{ID: seededUser.ID}
	err = user.DeleteByID(tokenDetails.AccessUUID, tokenDetails.RefreshUUID)
	if err != nil {
		return
	}

	assert.Equal(t, err, nil)
}

func TestFindUserById(t *testing.T) {
	_, err := seedOneUser()
	if err != nil {
		fmt.Printf("Error seedOneUser: %v", err)
	}

	user := User{ID: 1, Name: "pet"}
	receivedUser, err := user.FindByID()
	if err != nil {
		t.Errorf("User FindByID: %v", err)
		return
	}

	assert.Equal(t, user.ID, receivedUser.ID)
	assert.Equal(t, user.Name, receivedUser.Name)
}

func TestFindUserByEmail(t *testing.T) {
	_, err := seedOneUser()
	if err != nil {
		fmt.Printf("Error seedOneUser: %v", err)
	}

	user := User{Email: "devs@gmail.com", Name: "pet"}
	receivedUser, err := user.FindByEmail()
	if err != nil {
		t.Errorf("User FindByID: %v", err)
		return
	}

	assert.Equal(t, user.Email, receivedUser.Email)
	assert.Equal(t, user.Name, receivedUser.Name)
}

func TestFindAllUsers(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatalf("Error refreshUserTable: %v", err)
	}

	_, err = seedUsers()
	if err != nil {
		log.Fatalf("Error seedUsers: %v", err)
	}

	var user User
	users, err := user.FindAll()
	if err != nil {
		t.Errorf("User FindAll: %v", err)
		return
	}

	assert.Equal(t, len(*users), 2)
}
