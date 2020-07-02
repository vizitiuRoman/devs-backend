package controller_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	. "github.com/devs-backend/user-service/pkg/auth"
	. "github.com/devs-backend/user-service/pkg/controllers"
	"github.com/devs-backend/user-service/pkg/models"
	"github.com/gorilla/mux"
	"gopkg.in/stretchr/testify.v1/assert"
)

func TestGetUserByID(t *testing.T) {
	user, err := seedOneUser()
	if err != nil {
		fmt.Printf("Error seedOneUser: %v", err)
	}

	samples := []sample{
		{
			id:         int(user.ID),
			statusCode: 200,
			name:       user.Name,
			email:      user.Email,
		},
		{
			id:         10,
			statusCode: 404,
		},
	}

	for _, v := range samples {
		req, err := http.NewRequest("GET", "/users", nil)
		if err != nil {
			t.Errorf("Error NewRequest: %v", err)
		}

		req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(v.id)})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetUserByID)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, user.Name, responseMap["name"])
			assert.Equal(t, user.Email, responseMap["email"])
		}
	}
}

func TestGetUsers(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatalf("Error refreshUserTable: %v", err)
	}

	_, err = seedUsers()
	if err != nil {
		log.Fatalf("Error seedUsers: %v", err)
	}

	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Errorf("this is the error: %v", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetUsers)
	handler.ServeHTTP(rr, req)

	var users []models.User
	err = json.Unmarshal([]byte(rr.Body.String()), &users)
	if err != nil {
		log.Fatalf("Cannot convert to json: %v", err)
	}

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(users), 2)
}

func TestDeleteUser(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatalf("Error refreshUserTable: %v", err)
	}

	users, err := seedUsers()
	if err != nil {
		log.Fatalf("Error seedUsers: %v", err)
	}

	userID := users[0].ID
	token, err := CreateToken(userID)
	if err != nil {
		log.Fatalf("Error CreateToken: %v", err)
	}

	bearerToken := fmt.Sprintf("Bearer %v", token)
	samples := []sample{
		{
			statusCode:   422,
			token:        "",
			name:         "gopher",
			email:        "gopher@gmail.com",
			lastName:     "",
			errorMessage: "Unauthorized",
		},
		{
			uid:          strconv.FormatUint(uint64(userID), 10),
			statusCode:   200,
			token:        bearerToken,
			name:         "gopher",
			email:        "gopher@gmail.com",
			lastName:     "chitaica",
			errorMessage: "",
		},
		{
			uid:          strconv.FormatUint(uint64(userID), 10),
			statusCode:   401,
			token:        "bearerToken",
			errorMessage: "Unauthorized",
		},
		{
			uid:          strconv.FormatUint(uint64(userID), 10),
			statusCode:   500,
			token:        bearerToken,
			errorMessage: "Internal Server Error",
		},
	}

	for _, v := range samples {
		req, err := http.NewRequest("GET", "/users", nil)
		if err != nil {
			t.Errorf("This is the error: %v", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.uid})
		req.Header.Set("Authorization", v.token)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(DeleteUserByID)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 401 && v.errorMessage != "" {
			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestUpdateUser(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatalf("Error refreshUserTable: %v", err)
	}

	users, err := seedUsers()
	if err != nil {
		log.Fatalf("Error seedUsers: %v", err)
	}

	user := users[0]
	token, err := CreateToken(user.ID)
	if err != nil {
		log.Fatalf("Error CreateToken: %v", err)
	}

	bearerToken := fmt.Sprintf("Bearer %v", token)
	samples := []sample{
		{
			inputJSON:    `{"email": "gopher@gmail.com", "name": "", "lastName": "chitaica", "password": "devsmd"}`,
			statusCode:   422,
			token:        "",
			name:         "",
			email:        "gopher@gmail.com",
			lastName:     "gopher",
			errorMessage: "Unprocessable Entity",
		},
		{
			uid:          "1",
			inputJSON:    `{"email": "gopher@gmail.com", "name": "", "lastName": "chitaica", "password": "devsmd"}`,
			statusCode:   400,
			token:        "",
			name:         "",
			email:        "gopher@gmail.com",
			lastName:     "gopher",
			errorMessage: "Required Name",
		},
		{
			uid:          "1",
			inputJSON:    `{"email": "gopher@gmail.com", "name": "gopher", "lastName": "", "password": "devsmd"}`,
			statusCode:   400,
			token:        "",
			name:         "gopher",
			email:        "gopher@gmail.com",
			lastName:     "",
			errorMessage: "Required Last Name",
		},
		{
			uid:          "1",
			inputJSON:    `{"email": "devsmd", "name": "gopher", "lastName": "chitaica", "password": "devsmd"}`,
			statusCode:   400,
			token:        "",
			name:         "gopher",
			email:        "devsmd",
			lastName:     "chitaica",
			errorMessage: "Invalid Email",
		},
		{
			uid:          "1",
			inputJSON:    `{"email": "devsmd", "name": "gopher", "lastName": "chitaica", "password": "devsmd"}`,
			statusCode:   400,
			token:        "",
			name:         "gopher",
			email:        "devsmd",
			lastName:     "chitaica",
			errorMessage: "Invalid Email",
		},
		{
			uid:          "3",
			inputJSON:    `{"email": "gopher@gmail.com", "name": "gopher", "lastName": "chitaica", "password": "devsmd"}`,
			statusCode:   401,
			token:        bearerToken,
			name:         "gopher",
			email:        "gopher@gmail.com",
			lastName:     "chitaica",
			errorMessage: "Unauthorized",
		},
		{
			uid:          "1",
			inputJSON:    `{"email": "gopher@gmail.com", "name": "Updated", "lastName": "Good", "password": "devsmd"}`,
			statusCode:   200,
			token:        bearerToken,
			name:         "gopher",
			email:        "gopher@gmail.com",
			lastName:     "chitaica",
			errorMessage: "Unauthorized",
		},
	}

	for _, v := range samples {
		req, err := http.NewRequest("PUT", "/users", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("This is the error: %v", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.uid})
		req.Header.Set("Authorization", v.token)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(UpdateUser)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 401 && v.errorMessage != "" {
			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
