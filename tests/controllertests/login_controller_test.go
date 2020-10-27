package controllertests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/go-playground/assert.v1"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test the sign-in controller.
func TestSignIn(t *testing.T) {
	err := refreshTables()
	if err != nil {
		log.Fatalf("Failed to refresh tables: %v\n", err)
	}

	user, err := seedOneUser()
	if err != nil {
		fmt.Printf("Failed to seed user table: %v\n", err)
	}

	// sample user payloads
	mockData := []struct {
		email        string
		password     string
		errorMessage string
	}{
		{
			email:        user.Email,
			password:     user.Password,
			errorMessage: "",
		},
		{
			email:        user.Email,
			password:     "wrong password",
			errorMessage: "crypto/bcrypt: hashedPassword is not the hash of the given password",
		},
		{
			email:        "wrong email",
			password:     "wrong password",
			errorMessage: "record not found",
		},
	}

	// test the sample payloads
	for _, v := range mockData {
		token, err := server.AuthenticateCredentials(v.email, v.password)
		if err != nil {
			assert.Equal(t, err, errors.New(v.errorMessage))
		} else {
			assert.NotEqual(t, token, "")
		}
	}
}

// Test the login controller.
func TestLogin(t *testing.T) {
	err := refreshTables()
	if err != nil {
		log.Fatalf("Failed to refresh tables: %v\n", err)
	}

	// mock user
	user, err := seedOneUser()
	if err != nil {
		fmt.Printf("Failed to seed user table: %v\n", err)
	}

	// mock request JSON payloads
	correctCredentials := fmt.Sprintf(`{"email": "%v" , "password": "%v"}`, user.Email, user.Password)
	wrongPassword := fmt.Sprintf(`{"email": "%v" , "password": "%v"}`, user.Email, "wrong password")
	wrongCredentials := fmt.Sprintf(`{"email": "%v" , "password": "%v"}`, "wrong@email.com", "wrong password")
	invalidEmail := fmt.Sprintf(`{"email": "%v" , "password": "%v"}`, "invalid email", user.Password)
	missingEmail := fmt.Sprintf(`{"email": "%v" , "password": "%v"}`, "", user.Password)
	missingPassword := fmt.Sprintf(`{"email": "%v" , "password": "%v"}`, user.Email, "")

	// sample request payloads and responses
	samples := []struct {
		inputJSON    string
		statusCode   int
		email        string
		password     string
		errorMessage string
	}{
		{
			inputJSON:    correctCredentials,
			statusCode:   200,
			errorMessage: "",
		},
		{
			inputJSON:    wrongPassword,
			statusCode:   422,
			errorMessage: "incorrect password",
		},
		{
			inputJSON:    wrongCredentials,
			statusCode:   422,
			errorMessage: "incorrect details",
		},
		{
			inputJSON:    invalidEmail,
			statusCode:   422,
			errorMessage: "invalid email",
		},
		{
			inputJSON:    missingEmail,
			statusCode:   422,
			errorMessage: "email required",
		},
		{
			inputJSON:    missingPassword,
			statusCode:   422,
			errorMessage: "password required",
		},
	}

	// test the sample requests
	for _, v := range samples {
		// build the request
		req, err := http.NewRequest("POST", "/login", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("Error calling /login: %v\n", err)
		}
		// create response recorder and serve the request
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.Login)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, v.statusCode)
		// valid request tests
		if v.statusCode == 200 {
			assert.NotEqual(t, rr.Body.String(), "")
		}
		// invalid request tests
		if v.statusCode == 422 && v.errorMessage != "" {
			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				t.Errorf("Cannot convert to json: %v\n", err)
			}
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
