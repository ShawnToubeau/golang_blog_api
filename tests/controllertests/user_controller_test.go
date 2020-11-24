package controllertests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/shawntoubeau/golang_blog_api/api/models"
	"github.com/shawntoubeau/golang_blog_api/api/seed"
	"gopkg.in/go-playground/assert.v1"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestCreateUser(t *testing.T) {
	// seed test data
	_, _ = seed.Load(server.DB)
	// create new user
	newUser1 := seed.GenerateNewUser("Tony", "tony@tonyspizza.com", "pizza")
	newUser2 := seed.GenerateNewUser("Stevie", "stevie@hottopic.com", "pain")

	// mock request payloads
	validPayload := fmt.Sprintf(`{"nickname": "%v", "email": "%v" , "password": "%v"}`, newUser1.Nickname, newUser1.Email, newUser1.Password)
	emailTaken := fmt.Sprintf(`{"nickname": "%v", "email": "%v" , "password": "%v"}`, newUser2.Nickname, newUser1.Email, newUser2.Password)
	nicknameTaken := fmt.Sprintf(`{"nickname": "%v", "email": "%v" , "password": "%v"}`, newUser1.Nickname, newUser2.Email, newUser2.Password)
	emailMissing := fmt.Sprintf(`{"nickname": "%v", "email": "" , "password": "%v"}`, newUser2.Nickname, newUser2.Password)
	nicknameMissing := fmt.Sprintf(`{"nickname": "", "email": "%v" , "password": "%v"}`, newUser2.Email, newUser2.Password)
	passwordMissing := fmt.Sprintf(`{"nickname": "%v", "email": "%v" , "password": ""}`, newUser2.Nickname, newUser2.Email)

	// sample request payloads and responses
	samples := []struct {
		inputJSON    string
		statusCode   int
		nickname     string
		email        string
		errorMessage string
	}{
		// valid
		{
			validPayload,
			201,
			newUser1.Nickname,
			newUser1.Email,
			"",
		},
		// email taken
		{
			emailTaken,
			500,
			newUser2.Nickname,
			newUser1.Email,
			"email already taken",
		},
		// nickname taken
		{
			nicknameTaken,
			500,
			newUser1.Nickname,
			newUser2.Email,
			"nickname already taken",
		},
		// email missing
		{
			emailMissing,
			422,
			newUser2.Nickname,
			"",
			"email required",
		},
		// nickname missing
		{
			nicknameMissing,
			422,
			"",
			newUser2.Email,
			"nickname required",
		},
		// password missing
		{
			passwordMissing,
			422,
			newUser2.Nickname,
			newUser2.Email,
			"password required",
		},
	}

	// test sample requests
	for _, sample := range samples {
		// build the request
		req, err := http.NewRequest("PUT", "/user", bytes.NewBufferString(sample.inputJSON))
		if err != nil {
			t.Errorf("Error creating request: %v\n", err)
		}
		// create response recorder and serve the request
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.InsertUser)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to JSON: %v\n", err)
		}

		assert.Equal(t, rr.Code, sample.statusCode)
		// valid request tests
		if sample.statusCode == 201 {
			assert.Equal(t, responseMap["nickname"], sample.nickname)
			assert.Equal(t, responseMap["email"], sample.email)
		}
		// invalid request tests
		if sample.statusCode == 422 || sample.statusCode == 500 && sample.errorMessage != "" {
			assert.Equal(t, responseMap["error"], sample.errorMessage)
		}
	}
}

func TestFetchUsers(t *testing.T) {
	// seed test data
	users, _ := seed.Load(server.DB)
	// create request
	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Errorf("Failed to create request: %v\n", err)
	}
	// create request recorder
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.FetchAllUsers)
	// serve request
	handler.ServeHTTP(rr, req)
	// create user array and process response
	var fetchedUsers []models.User
	err = json.Unmarshal([]byte(rr.Body.String()), &fetchedUsers)
	if err != nil {
		log.Fatalf("Cannot convert to JSON: %v\n", err)
	}

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(fetchedUsers), len(users))
}

func TestFetchUserById(t *testing.T) {
	// seed test data
	users, _ := seed.Load(server.DB)
	user := users[0]
	user.Password = seed.MockUser1.Password

	// sample request payloads and responses
	userSample := []struct {
		id         string
		statusCode int
		nickname   string
		email      string
		password   string
	}{
		{
			strconv.Itoa(int(user.ID)),
			200,
			user.Nickname,
			user.Email,
			user.Password,
		},
		{
			"unknown",
			400,
			"",
			"",
			"",
		},
	}

	// test the sample requests
	for _, sample := range userSample {
		// create request
		req, err := http.NewRequest("GET", "/users", nil)
		if err != nil {
			t.Errorf("Failed to create request: %v\n", err)
		}
		// set request params, create response recorder, and serve the request
		req = mux.SetURLVars(req, map[string]string{"id": sample.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.FetchUserByID)
		handler.ServeHTTP(rr, req)
		// store the response
		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to JSON: %v\n", err)
		}

		assert.Equal(t, rr.Code, sample.statusCode)
		// valid request tests
		if sample.statusCode == 200 {
			assert.Equal(t, sample.nickname, responseMap["nickname"])
			assert.Equal(t, sample.email, responseMap["email"])
		}
	}
}

func TestUpdateUser(t *testing.T) {
	// seed test data
	users, _ := seed.Load(server.DB)
	// retrieve first user
	user := users[0]
	user.Password = seed.MockUser1.Password
	// login user to retrieve auth token
	token, err := server.AuthenticateCredentials(user.Email, user.Password)
	if err != nil {
		log.Fatalf("Failed to login user: %v\n", err)
	}
	// format token
	tokenString := fmt.Sprintf("Bearer %v", token)

	// mock request JSON payloads
	validPayload := fmt.Sprintf(`{"nickname": "%v", "email": "%v", "password": "%v"}`, user.Nickname, user.Email, user.Password)
	missingPassword := fmt.Sprintf(`{"nickname": "%v", "email": "%v", "password": ""}`, user.Nickname, user.Email)
	emailTaken := fmt.Sprintf(`{"nickname": "%v", "email": "%v", "password": "%v"}`, user.Nickname, users[1].Email, user.Password)
	nicknameTaken := fmt.Sprintf(`{"nickname": "%v", "email": "%v", "password": "%v"}`, users[1].Nickname, user.Email, user.Password)
	invalidEmail := fmt.Sprintf(`{"nickname": "%v", "email": "%v", "password": "%v"}`, user.Nickname, "invalid email", user.Password)
	missingEmail := fmt.Sprintf(`{"nickname": "%v", "email": "%v", "password": "%v"}`, user.Nickname, "", user.Password)
	missingNickname := fmt.Sprintf(`{"nickname": "%v", "email": "%v", "password": "%v"}`, "", user.Email, user.Password)

	// sample request payloads and responses
	samples := []struct {
		id             string
		updateJSON     string
		statusCode     int
		updateNickname string
		updateEmail    string
		tokenGiven     string
		errorMessage   string
	}{
		// Valid
		{
			strconv.Itoa(int(user.ID)),
			validPayload,
			200,
			user.Nickname,
			user.Email,
			tokenString,
			"",
		},
		// Empty password field
		{
			strconv.Itoa(int(user.ID)),
			missingPassword,
			422,
			user.Nickname,
			user.Email,
			tokenString,
			"password required",
		},
		// No auth token
		{
			strconv.Itoa(int(user.ID)),
			validPayload,
			401,
			user.Nickname,
			user.Email,
			"",
			"Unauthorized",
		},
		// Wrong auth token
		{
			strconv.Itoa(int(user.ID)),
			validPayload,
			401,
			user.Nickname,
			user.Email,
			"wrong token",
			"Unauthorized",
		},
		// Email taken
		{
			strconv.Itoa(int(user.ID)),
			emailTaken,
			500,
			user.Nickname,
			users[1].Email,
			tokenString,
			"email already taken",
		},
		// Nickname taken
		{
			strconv.Itoa(int(user.ID)),
			nicknameTaken,
			500,
			users[1].Nickname,
			user.Email,
			tokenString,
			"nickname already taken",
		},
		// Email invalid
		{
			strconv.Itoa(int(user.ID)),
			invalidEmail,
			422,
			user.Nickname,
			"invalid email",
			tokenString,
			"invalid email",
		},
		// Email field empty
		{
			strconv.Itoa(int(user.ID)),
			missingEmail,
			422,
			user.Nickname,
			"",
			tokenString,
			"email required",
		},
		// Nickname field empty
		{
			strconv.Itoa(int(user.ID)),
			missingNickname,
			422,
			"",
			user.Email,
			tokenString,
			"nickname required",
		},
		// No ID
		{
			"",
			validPayload,
			400,
			user.Nickname,
			user.Email,
			tokenString,
			"no ID",
		},
		// Using other user's token
		{
			strconv.Itoa(2),
			validPayload,
			401,
			user.Nickname,
			user.Email,
			tokenString,
			"Unauthorized",
		},
	}

	// test the sample requests
	for _, sample := range samples {
		// create request
		req, err := http.NewRequest("POST", "/users", bytes.NewBufferString(sample.updateJSON))
		if err != nil {
			t.Errorf("Failed to create request: %v\n", err)
		}
		// set request params
		req = mux.SetURLVars(req, map[string]string{"id": sample.id})
		req.Header.Set("Authorization", sample.tokenGiven)
		// create response recorder
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdateUserById)
		// serve the request
		handler.ServeHTTP(rr, req)
		// store the response
		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to JSON: %v\n", err)
		}

		assert.Equal(t, rr.Code, sample.statusCode)
		// valid request tests
		if sample.statusCode == 200 {
			assert.Equal(t, responseMap["nickname"], sample.updateNickname)
			assert.Equal(t, responseMap["email"], sample.updateEmail)
		}
		// invalid request tests
		if sample.statusCode == 401 || sample.statusCode == 422 || sample.statusCode == 500 && sample.errorMessage != "" {
			assert.Equal(t, responseMap["error"], sample.errorMessage)
		}
	}
}

func TestDeleteUser(t *testing.T) {
	// seed test data
	users, posts := seed.Load(server.DB)
	// retrieve first user
	user := users[0]
	user.Password = seed.MockUser1.Password
	// login the user to get their auth token
	token, err := server.AuthenticateCredentials(user.Email, user.Password)
	if err != nil {
		log.Fatalf("Failed to login user: %v\n", err)
	}
	// construct token string
	tokenString := fmt.Sprintf("Bearer %v", token)

	// sample request payloads and responses
	samples := []struct {
		id           string
		tokenGiven   string
		stateCode    int
		errorMessage string
	}{
		// User contains linked posts
		{
			strconv.Itoa(int(user.ID)),
			tokenString,
			500,
			"ERROR: update or delete on table \"users\" violates foreign key constraint \"fk_posts_author\" on table \"posts\" (SQLSTATE 23503)",
		},
		// Valid
		{
			strconv.Itoa(int(user.ID)),
			tokenString,
			204,
			"",
		},
		// Missing token string
		{
			strconv.Itoa(int(user.ID)),
			"",
			401,
			"Unauthorized",
		},
		// Incorrect token string
		{
			strconv.Itoa(int(user.ID)),
			"Incorrect token",
			401,
			"Unauthorized",
		},
		// Missing user ID
		{
			"",
			tokenString,
			400,
			"",
		},
		// Wrong user ID
		{
			strconv.Itoa(int(2)),
			tokenString,
			401,
			"Unauthorized",
		},
	}

	// test each sample request payload
	for i, sample := range samples {
		// build the request
		req, err := http.NewRequest("DELETE", "/users", nil)
		if err != nil {
			t.Errorf("Failed to create request: %v\n", err)
		}
		// set request variables and create response recorder
		req = mux.SetURLVars(req, map[string]string{"id": sample.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeleteUserById)
		// set token header
		req.Header.Set("Authorization", sample.tokenGiven)
		// serve the request
		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, sample.stateCode)
		// failed request tests
		if sample.stateCode != 204 && sample.errorMessage != "" {
			// create response map
			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				t.Errorf("Cannot convert to json: %v\n", err)
			}
			assert.Equal(t, responseMap["error"], sample.errorMessage)

			// delete all posts so there will no foreign key violations
			// when trying to delete the user again
			if i == 0 {
				for _, post := range posts {
					post.DeletePostById(server.DB, post.ID, post.AuthorID)
				}
			}
		}
	}
}
