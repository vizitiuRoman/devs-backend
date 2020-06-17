package controllertests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/devs-backend/user-service/pkg/controllers"
	"gopkg.in/stretchr/testify.v1/assert"
)

func TestRegisterUser(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	samples := []struct {
		inputJSON    string
		statusCode   int
		name         string
		email        string
		lastName     string
		token        string
		errorMessage string
	}{
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
			t.Errorf("this is the error: %v", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(Register)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.NotEmpty(t, responseMap["token"])
			assert.Equal(t, responseMap["name"], v.name)
			assert.Equal(t, responseMap["lastName"], v.lastName)
		}
		if v.statusCode == 400 || v.statusCode == 409 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
