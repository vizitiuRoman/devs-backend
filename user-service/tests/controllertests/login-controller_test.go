package controllertests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/devs-backend/user-service/pkg/auth"
	. "github.com/devs-backend/user-service/pkg/controllers"
	"gopkg.in/stretchr/testify.v1/assert"
)

type sample struct {
	inputJSON    string
	statusCode   int
	name         string
	email        string
	lastName     string
	token        string
	errorMessage string
}

func TestRegisterUser(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatalf("Error refreshing user table %v", err)
	}

	samples := []sample{
		{
			inputJSON:    `{"email": "devsmd@gmail.com", "name": "", "lastName": "Devs", "password": "devsmd"}`,
			statusCode:   400,
			name:         "",
			email:        "devsmd@gmail.com",
			lastName:     "Devs",
			errorMessage: "Required Name",
		},
		{
			inputJSON:    `{"email": "devsmd@gmail.com", "name": "Devs", "lastName": "", "password": "devsmd"}`,
			statusCode:   400,
			name:         "Devs",
			email:        "devsmd@gmail.com",
			lastName:     "",
			errorMessage: "Required Last Name",
		},
		{
			inputJSON:    `{"email": "devsmd@gmail.com", "name": "Devs", "lastName": "Devs", "password": ""}`,
			statusCode:   400,
			name:         "Devs",
			email:        "devsmd@gmail.com",
			lastName:     "Devs",
			errorMessage: "Required Password",
		},
		{
			inputJSON:    `{"email": "devsmd", "name": "Devs", "lastName": "Devs", "password": "devsmd"}`,
			statusCode:   400,
			name:         "Devs",
			email:        "devsmd",
			lastName:     "Devs",
			errorMessage: "Invalid Email",
		},
		{
			inputJSON:    `{"email": "devsmd@gmail.com", "name": "Devs", "lastName": "Devs", "password": "devsmd"}`,
			statusCode:   201,
			name:         "Devs",
			email:        "devsmd@gmail.com",
			lastName:     "Devs",
			errorMessage: "",
		},
		{
			inputJSON:    `{"email": "devsmd@gmail.com", "name": "Devs", "lastName": "Devs", "password": "devsmd"}`,
			statusCode:   409,
			name:         "Devs",
			email:        "devsmd@gmail.com",
			lastName:     "Devs",
			errorMessage: "Conflict",
		},
	}

	for _, v := range samples {
		req, err := http.NewRequest("POST", "/register", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("NewRequest: %v", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(Register)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			bearerToken := "Bearer " + responseMap["token"].(string)
			req.Header.Set("Authorization", bearerToken)

			assert.NotEmpty(t, responseMap["token"])
			assert.Equal(t, responseMap["name"], v.name)
			assert.Equal(t, responseMap["lastName"], v.lastName)

			token, err := ExtractTokenMetadata(req)
			if err != nil {
				t.Errorf("ExtractTokenMetadata: %v", err)
			}
			_, err = FetchToken(&AccessDetails{
				AccessUUID: token.AccessUUID,
				UserID:     token.UserID,
			})
			if err != nil {
				t.Errorf("FetchToken: %v", err)
			}
		}
		if v.statusCode >= 400 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestLoginUser(t *testing.T) {
	_, err := seedOneUser()
	if err != nil {
		fmt.Printf("seedOneUser: %v", err)
	}

	samples := []sample{
		{
			inputJSON:    `{"email": "devsmd@gmail.com", "password": "devsm"}`,
			statusCode:   400,
			name:         "",
			email:        "devsmd@gmail.com",
			lastName:     "Devs",
			errorMessage: "Invalid Email Or Password",
		},
		{
			inputJSON:    `{"email": "devsmd@gmail.co", "password": "devsmd"}`,
			statusCode:   400,
			name:         "Devs",
			email:        "devsmd@gmail.com",
			lastName:     "Devs",
			errorMessage: "Invalid Email Or Password",
		},
		{
			inputJSON:    `{"email": "devs@gmail.com", "password": "password"}`,
			statusCode:   200,
			name:         "pet",
			lastName:     "pets",
			errorMessage: "",
		},
	}

	for _, v := range samples {
		req, err := http.NewRequest("POST", "/login", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("NewRequest: %v", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(Login)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 || v.statusCode == 200 {
			bearerToken := "Bearer " + responseMap["token"].(string)
			req.Header.Set("Authorization", bearerToken)

			assert.NotEmpty(t, responseMap["token"])
			assert.Equal(t, responseMap["name"], v.name)
			assert.Equal(t, responseMap["lastName"], v.lastName)

			token, err := ExtractTokenMetadata(req)
			if err != nil {
				t.Errorf("ExtractTokenMetadata: %v", err)
			}
			_, err = FetchToken(&AccessDetails{
				AccessUUID: token.AccessUUID,
				UserID:     token.UserID,
			})
			if err != nil {
				t.Errorf("FetchToken: %v", err)
			}
		}

		if v.statusCode >= 400 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestLogoutUser(t *testing.T) {
	_, err := seedOneUser()
	if err != nil {
		fmt.Printf("seedOneUser: %v", err)
	}

	v := sample{
		inputJSON:    `{"email": "devs@gmail.com", "password": "password"}`,
		statusCode:   200,
		name:         "pet",
		lastName:     "pets",
		errorMessage: "",
	}

	req, err := http.NewRequest("POST", "/login", bytes.NewBufferString(v.inputJSON))
	if err != nil {
		t.Errorf("NewRequest: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Login)
	handler.ServeHTTP(rr, req)

	responseMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
	if err != nil {
		t.Errorf("Cannot convert to json: %v", err)
	}

	assert.Equal(t, rr.Code, v.statusCode)
	if v.statusCode == 200 {
		rq, err := http.NewRequest("POST", "/logout", bytes.NewBufferString(""))
		if err != nil {
			t.Errorf("NewRequest: %v", err)
		}

		bearerToken := "Bearer " + responseMap["token"].(string)
		rq.Header.Set("Authorization", bearerToken)

		r := httptest.NewRecorder()
		handler := http.HandlerFunc(Logout)
		handler.ServeHTTP(r, rq)

		assert.Equal(t, r.Code, v.statusCode)
	}
}
