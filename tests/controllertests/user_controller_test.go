package controllertests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"github.com/gorilla/mux"
	"github.com/shawntoubeau/golang_blog_api/api/models"
	"gopkg.in/go-playground/assert.v1"
)

func TestCreateUser(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	samples := []struct{
		inputJSON string
		statusCode int
		nickname string
		email string
		errorMessage string
	}{
		{
			`{"nickname": "Shawn", "email": "shawn@aol.com", "password": "123"}`,
			201,
			"shawn",
			"shawn@aol.com",
			"",
		},
		{
			`{"nickname": "Aria", "email": "shawn@aol.com", "password": "321"}`,
			500,
			nil,
			nil,
			"Email Already Taken",
		},
		{
			`{"nickname": "Shawn", "email": "aria@aol.com", "password": "321"}`,
			500,
			nil,
			nil,
			"Nickname Already Taken",
		},
		{
			`{"nickname": "", "email": "aria@aol.com", "password": "321"}`,
			422,
			nil,
			nil,
			"Invalid Email",
		},
		{
			`{"nickname": "", "email": "aria@aol.com", "password": "321"}`,
			422,
			nil,
			nil,
			"Required Nickname",
		},
		{
			`{"nickname": "Aria", "email": "", "password": "321"}`,
			422,
			nil,
			nil,
			"Required Email",
		},
		{
			`{"nickname": "Aria", "email": "aria@aol.com", "password": ""}`,
			422,
			nil,
			nil,
			"Required Password",
		},
	}

	for _, v := range samples {
		req, err := http.NewRequest("POST", "/user", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("Error creating request: %v\n", err)
		}
		// Records server responses
		rr := httptest.NewRecorder()
		// Request handler
		handler := http.HandlerFunc(server.CreateUser)
		// Serve request
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to JSON: %v\n")
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["nickname"], v.nickname)
			assert.Equal(t, responseMap["email"], v.email)
		}
		if v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetUsers(t *testing.T) {
	// Refresh table
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	// Seed table
	_, err = seedUsers()
	if err != nil {
		log.Fatal(err)
	}
	// Create request
	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Errorf("Failed to form request: %v\n", err)
	}
	// Create request recorder and serve
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetUsers)
	handler.ServeHTTP(rr, req)
	// Create user array and process response
	var users []models.User
	err = json.Unmarshal([]byte(rr.Body.String()), &users)
	if err != nil {
		log.Fatalf("Cannot convert to JSON: %v\n", err)
	}

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(users), 2)
}

func TestGetUserById(t *testing.T) {
	// Refresh table
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	// Seed v
	user, err := seedOneUser()
	if err != nil {
		log.Fatal(err)
	}
	userSample := []struct {
		id string
		statusCode int
		nickname string
		email string
		password string
	}{
		{

			strconv.Itoa(int(user.ID)),
			200,
			user.Nickname,
			user.Email,
			nil,
		},
		{
			"unknown",
			400,
			nil,
			nil,
			nil,
		},
	}

	for _, v := range userSample {
		// Create request
		req, err := http.NewRequest("GET", "/users", nil)
		if err != nil {
			t.Errorf("Failed to create request: %v\n", err)
		}
		// Set request params and test request
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetUserByID)
		handler.ServeHTTP(rr, req)
		// Process response
		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to JSON: %v\n", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, user.Nickname, responseMap["nickname"])
			assert.Equal(t, user.Email, responseMap["email"])
		}
	}
}

func TestUpdateUser(t *testing.T) {
	var AuthEmail, AuthPassword string
	var AuthID uint32

	// Refresh user table
	err := refreshUserTable();
	if err != nil {
		log.Fatalf("Failed to refresh user table: %v\n", err)
	}
	// Seed one user
	users, err := seedUsers()
	if err != nil {
		log.Fatalf("Failed to seed user table: %v\n", err)
	}
	// Retrieve first user
	for _, user := range users {
		if user.ID == 2 {
			continue
		}
		AuthID = user.ID
		AuthEmail = user.Email
		AuthPassword = user.Password
	}
	// Login user to retrieve auth token
	token, err := server.SignIn(AuthEmail, AuthPassword)
	if err != nil {
		log.Fatalf("Failed to login user: %v\n", err)
	}
	// Format token
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		id string
		updateJSON string
		statusCode int
		updateNickname string
		updateEmail string
		tokenGiven string
		errorMessage string
	} {
		// OK
		{
			strconv.Itoa(int(AuthID)),
			`{"nickname": "Shawn", "email": "shawn@aol.com", "password": "123"}`,
			200,
			"Swan",
			"swan@aol.com",
			tokenString,
			nil,
		},
		// Empty password field
		{
			strconv.Itoa(int(AuthID)),
			`{"nickname": "Shawn", "email": "shawn@aol.com", "password": ""}`,
			422,
			"Swan",
			"swan@aol.com",
			tokenString,
			"Required Password",
		},
		// No auth token
		{
			strconv.Itoa(int(AuthID)),
			`{"nickname": "Shawn", "email": "shawn@aol.com", "password": "123"}`,
			200,
			"Swan",
			"swan@aol.com",
			"",
			"Unauthorized",
		},
		// Wrong auth token
		{
			strconv.Itoa(int(AuthID)),
			`{"nickname": "Shawn", "email": "shawn@aol.com", "password": "123"}`,
			401,
			"Swan",
			"swan@aol.com",
			"wrong token",
			"Unauthorized",
		},
		// Email taken
		{
			strconv.Itoa(int(AuthID)),
			`{"nickname": "Shawn", "email": "shawn@aol.com", "password": "123"}`,
			500,
			"Swan",
			"aria@aol.com",
			tokenString,
			"Email Already Taken",
		},
		// Nickname taken
		{
			strconv.Itoa(int(AuthID)),
			`{"nickname": "Shawn", "email": "shawn@aol.com", "password": "123"}`,
			200,
			"Aria",
			"swan@aol.com",
			tokenString,
			"Nickname Already Taken",
		},
		// Email invalid
		{
			strconv.Itoa(int(AuthID)),
			`{"nickname": "Shawn", "email": "shawn@aol.com", "password": "123"}`,
			422,
			"Swan",
			nil,
			tokenString,
			"Invalid Email",
		},
		// Email field empty
		{
			strconv.Itoa(int(AuthID)),
			`{"nickname": "Shawn", "email": "", "password": "123"}`,
			422,
			nil,
			nil,
			tokenString,
			"Required Email",
		},
		// Nickname field empty
		{
			strconv.Itoa(int(AuthID)),
			`{"nickname": "", "email": "shawn@aol.com", "password": "123"}`,
			422,
			nil,
			nil,
			tokenString,
			"Required Nickname",
		},
		// No ID
		{
			nil,
			nil,
			400,
			nil,
			nil,
			tokenString,
			nil,
		},
		// Using other user's token
		{
			strconv.Itoa(int(2)),
			`{"nickname": "Shawn", "email": "shawn@aol.com", "password": "123"}`,
			401,
			nil,
			nil,
			tokenString,
			"Unauthorized",
		},

	}
}