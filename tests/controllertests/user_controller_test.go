package controllertests
//
//import (
//	"bytes"
//	"encoding/json"
//	"fmt"
//	"github.com/gorilla/mux"
//	"github.com/shawntoubeau/golang_blog_api/api/models"
//	"gopkg.in/go-playground/assert.v1"
//	"log"
//	"net/http"
//	"net/http/httptest"
//	"strconv"
//	"testing"
//)
//
//func TestCreateUser(t *testing.T) {
//	err := refreshTables()
//	if err != nil {
//		log.Fatalf("Failed to refresh tables: %v\n", err)
//	}
//
//	// mock users
//	user1 := MockUser1
//	user2 := MockUser2
//
//	// mock request payloads
//	user1InputJSON := fmt.Sprintf(`{"nickname": "%v", "email": "%v" , "password": "%v"}`, user1.Nickname, user1.Email, user1.Password)
//	emailAlreadyTaken := fmt.Sprintf(`{"nickname": "%v", "email": "%v" , "password": "%v"}`, user2.Nickname, user1.Email, user2.Password)
//	nicknameAlreadyTaken := fmt.Sprintf(`{"nickname": "%v", "email": "%v" , "password": "%v"}`, user1.Nickname, user2.Email, user2.Password)
//	emailRequiredTaken := fmt.Sprintf(`{"nickname": "%v", "email": "" , "password": "%v"}`, user2.Nickname, user2.Password)
//	nicknameRequiredTaken := fmt.Sprintf(`{"nickname": "", "email": "%v" , "password": "%v"}`, user2.Email, user2.Password)
//	passwordMissing := fmt.Sprintf(`{"nickname": "%v", "email": "%v" , "password": ""}`, user2.Nickname, user2.Email)
//
//	// sample request payloads and responses
//	samples := []struct {
//		inputJSON    string
//		statusCode   int
//		nickname     string
//		email        string
//		errorMessage string
//	}{
//		{
//			user1InputJSON,
//			201,
//			user1.Nickname,
//			user1.Email,
//			"",
//		},
//		{
//			emailAlreadyTaken,
//			500,
//			user2.Nickname,
//			user1.Email,
//			"email already taken",
//		},
//		{
//			nicknameAlreadyTaken,
//			500,
//			user1.Nickname,
//			user2.Email,
//			"nickname already taken",
//		},
//		{
//			emailRequiredTaken,
//			422,
//			user2.Nickname,
//			"",
//			"email required",
//		},
//		{
//			nicknameRequiredTaken,
//			422,
//			"",
//			user2.Email,
//			"nickname required",
//		},
//		{
//			passwordMissing,
//			422,
//			user2.Nickname,
//			user2.Email,
//			"password required",
//		},
//	}
//
//	// test sample requests
//	for _, sample := range samples {
//		// build the request
//		req, err := http.NewRequest("PUT", "/user", bytes.NewBufferString(sample.inputJSON))
//		if err != nil {
//			t.Errorf("Error creating request: %v\n", err)
//		}
//		// create response recorder and serve the request
//		rr := httptest.NewRecorder()
//		handler := http.HandlerFunc(server.InsertUser)
//		handler.ServeHTTP(rr, req)
//
//		responseMap := make(map[string]interface{})
//		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
//		if err != nil {
//			fmt.Printf("Cannot convert to JSON: %v\n", err)
//		}
//
//		assert.Equal(t, rr.Code, sample.statusCode)
//		// valid request tests
//		if sample.statusCode == 201 {
//			assert.Equal(t, responseMap["nickname"], sample.nickname)
//			assert.Equal(t, responseMap["email"], sample.email)
//		}
//		// invalid request tests
//		if sample.statusCode == 422 || sample.statusCode == 500 && sample.errorMessage != "" {
//			assert.Equal(t, responseMap["error"], sample.errorMessage)
//		}
//	}
//}
//
//func TestFetchUsers(t *testing.T) {
//	// refresh table
//	err := refreshTables()
//	if err != nil {
//		log.Fatal(err)
//	}
//	// seed table
//	mockUsers, err := seedUsers()
//	if err != nil {
//		log.Fatal(err)
//	}
//	// create request
//	req, err := http.NewRequest("GET", "/users", nil)
//	if err != nil {
//		t.Errorf("Failed to form request: %v\n", err)
//	}
//	// create request recorder and serve
//	rr := httptest.NewRecorder()
//	handler := http.HandlerFunc(server.FetchAllUsers)
//	handler.ServeHTTP(rr, req)
//	// create user array and process response
//	var users []models.User
//	err = json.Unmarshal([]byte(rr.Body.String()), &users)
//	if err != nil {
//		log.Fatalf("Cannot convert to JSON: %v\n", err)
//	}
//
//	assert.Equal(t, rr.Code, http.StatusOK)
//	assert.Equal(t, len(users), len(mockUsers))
//}
//
//func TestFetchUserById(t *testing.T) {
//	// refresh table
//	err := refreshTables()
//	if err != nil {
//		log.Fatal(err)
//	}
//	// seed user
//	user, err := seedOneUser()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// sample request payloads and responses
//	userSample := []struct {
//		id         string
//		statusCode int
//		nickname   string
//		email      string
//		password   string
//	}{
//		{
//
//			strconv.Itoa(int(user.ID)),
//			200,
//			user.Nickname,
//			user.Email,
//			MockUser1.Password,
//		},
//		{
//			"unknown",
//			400,
//			"",
//			"",
//			"",
//		},
//	}
//
//	// test the sample requests
//	for _, sample := range userSample {
//		// create request
//		req, err := http.NewRequest("GET", "/users", nil)
//		if err != nil {
//			t.Errorf("Failed to create request: %v\n", err)
//		}
//		// set request params, create response recorder, and serve the request
//		req = mux.SetURLVars(req, map[string]string{"id": sample.id})
//		rr := httptest.NewRecorder()
//		handler := http.HandlerFunc(server.FetchUserByID)
//		handler.ServeHTTP(rr, req)
//		// store the response
//		responseMap := make(map[string]interface{})
//		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
//		if err != nil {
//			log.Fatalf("Cannot convert to JSON: %v\n", err)
//		}
//
//		assert.Equal(t, rr.Code, sample.statusCode)
//		// valid request tests
//		if sample.statusCode == 200 {
//			assert.Equal(t, sample.nickname, responseMap["nickname"])
//			assert.Equal(t, sample.email, responseMap["email"])
//		}
//	}
//}
//
//func TestUpdateUser(t *testing.T) {
//	var AuthEmail, AuthPassword string
//	var AuthID uint32
//	// refresh user table
//	err := refreshTables()
//	if err != nil {
//		log.Fatalf("Failed to refresh user table: %v\n", err)
//	}
//	// seed user
//	users, err := seedUsers()
//	if err != nil {
//		log.Fatalf("Failed to seed user table: %v\n", err)
//	}
//	user := users[0]
//	userPassword := MockUser1.Password
//
//	// retrieve first user
//	AuthID = user.ID
//	AuthEmail = user.Email
//	AuthPassword = MockUser1.Password
//	// login user to retrieve auth token
//	token, err := server.AuthenticateCredentials(AuthEmail, AuthPassword)
//	if err != nil {
//		log.Fatalf("Failed to login user: %v\n", err)
//	}
//	// format token
//	tokenString := fmt.Sprintf("Bearer %v", token)
//
//	// mock request JSON payloads
//	validPayload := fmt.Sprintf(`{"nickname": "%v", "email": "%v", "password": "%v"}`, user.Nickname, user.Email, userPassword)
//	missingPassword := fmt.Sprintf(`{"nickname": "%v", "email": "%v", "password": ""}`, user.Nickname, user.Email)
//	emailTaken := fmt.Sprintf(`{"nickname": "%v", "email": "%v", "password": "%v"}`, user.Nickname, users[1].Email, userPassword)
//	nicknameTaken := fmt.Sprintf(`{"nickname": "%v", "email": "%v", "password": "%v"}`, users[1].Nickname, user.Email, userPassword)
//	invalidEmail := fmt.Sprintf(`{"nickname": "%v", "email": "%v", "password": "%v"}`, user.Nickname, "invalid email", userPassword)
//	missingEmail := fmt.Sprintf(`{"nickname": "%v", "email": "%v", "password": "%v"}`, user.Nickname, "", userPassword)
//	missingNickname := fmt.Sprintf(`{"nickname": "%v", "email": "%v", "password": "%v"}`, "", user.Email, userPassword)
//
//	// sample request payloads and responses
//	samples := []struct {
//		id             string
//		updateJSON     string
//		statusCode     int
//		updateNickname string
//		updateEmail    string
//		tokenGiven     string
//		errorMessage   string
//	}{
//		// Valid
//		{
//			strconv.Itoa(int(AuthID)),
//			validPayload,
//			200,
//			user.Nickname,
//			user.Email,
//			tokenString,
//			"",
//		},
//		// Empty password field
//		{
//			strconv.Itoa(int(AuthID)),
//			missingPassword,
//			422,
//			user.Nickname,
//			user.Email,
//			tokenString,
//			"password required",
//		},
//		// No auth token
//		{
//			strconv.Itoa(int(AuthID)),
//			validPayload,
//			401,
//			user.Nickname,
//			user.Email,
//			"",
//			"Unauthorized",
//		},
//		// Wrong auth token
//		{
//			strconv.Itoa(int(AuthID)),
//			validPayload,
//			401,
//			user.Nickname,
//			user.Email,
//			"wrong token",
//			"Unauthorized",
//		},
//		// Email taken
//		{
//			strconv.Itoa(int(AuthID)),
//			emailTaken,
//			500,
//			user.Nickname,
//			users[1].Email,
//			tokenString,
//			"email already taken",
//		},
//		// Nickname taken
//		{
//			strconv.Itoa(int(AuthID)),
//			nicknameTaken,
//			500,
//			users[1].Nickname,
//			user.Email,
//			tokenString,
//			"nickname already taken",
//		},
//		// Email invalid
//		{
//			strconv.Itoa(int(AuthID)),
//			invalidEmail,
//			422,
//			user.Nickname,
//			"invalid email",
//			tokenString,
//			"invalid email",
//		},
//		// Email field empty
//		{
//			strconv.Itoa(int(AuthID)),
//			missingEmail,
//			422,
//			user.Nickname,
//			"",
//			tokenString,
//			"email required",
//		},
//		// Nickname field empty
//		{
//			strconv.Itoa(int(AuthID)),
//			missingNickname,
//			422,
//			"",
//			user.Email,
//			tokenString,
//			"nickname required",
//		},
//		// No ID
//		{
//			"",
//			validPayload,
//			400,
//			user.Nickname,
//			user.Email,
//			tokenString,
//			"no ID",
//		},
//		// Using other user's token
//		{
//			strconv.Itoa(2),
//			validPayload,
//			401,
//			user.Nickname,
//			user.Email,
//			tokenString,
//			"Unauthorized",
//		},
//	}
//
//	// test the sample requests
//	for _, sample := range samples {
//		// create request
//		req, err := http.NewRequest("POST", "/users", bytes.NewBufferString(sample.updateJSON))
//		if err != nil {
//			t.Errorf("Failed to create request: %v\n", err)
//		}
//		// set request params
//		req = mux.SetURLVars(req, map[string]string{"id": sample.id})
//		req.Header.Set("Authorization", sample.tokenGiven)
//		// create response recorder
//		rr := httptest.NewRecorder()
//		handler := http.HandlerFunc(server.UpdateUserById)
//		// serve the request
//		handler.ServeHTTP(rr, req)
//		// store the response
//		responseMap := make(map[string]interface{})
//		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
//		if err != nil {
//			t.Errorf("Cannot convert to JSON: %v\n", err)
//		}
//
//		assert.Equal(t, rr.Code, sample.statusCode)
//		// valid request tests
//		if sample.statusCode == 200 {
//			assert.Equal(t, responseMap["nickname"], sample.updateNickname)
//			assert.Equal(t, responseMap["email"], sample.updateEmail)
//		}
//		// invalid request tests
//		if sample.statusCode == 401 || sample.statusCode == 422 || sample.statusCode == 500 && sample.errorMessage != "" {
//			assert.Equal(t, responseMap["error"], sample.errorMessage)
//		}
//	}
//}
//
//func TestDeleteUser(t *testing.T) {
//	var AuthEmail, AuthPassword string
//	var AuthId uint32
//	// refresh user table
//	err := refreshTables()
//	if err != nil {
//		log.Fatalf("Failed to refresh user table: %v\n", err)
//	}
//	// seed users
//	users, err := seedUsers()
//	if err != nil {
//		log.Fatalf("Failed to seed user table: %v\n", err)
//	}
//	// get first users credentials
//	AuthId = users[0].ID
//	AuthEmail = users[0].Email
//	AuthPassword = MockUser1.Password
//	// login the user to get their auth token
//	token, err := server.AuthenticateCredentials(AuthEmail, AuthPassword)
//	if err != nil {
//		log.Fatalf("Failed to login user: %v\n", err)
//	}
//	// construct token string
//	tokenString := fmt.Sprintf("Bearer %v", token)
//
//	// sample request payloads and responses
//	userSample := []struct {
//		id           string
//		tokenGiven   string
//		stateCode    int
//		errorMessage string
//	}{
//		// Valid
//		{
//			strconv.Itoa(int(AuthId)),
//			tokenString,
//			204,
//			"",
//		},
//		// Missing token string
//		{
//			strconv.Itoa(int(AuthId)),
//			"",
//			401,
//			"Unauthorized",
//		},
//		// Incorrect token string
//		{
//			strconv.Itoa(int(AuthId)),
//			"Incorrect token",
//			401,
//			"Unauthorized",
//		},
//		// Missing user ID
//		{
//			"",
//			tokenString,
//			400,
//			"",
//		},
//		// Wrong user ID
//		{
//			strconv.Itoa(int(2)),
//			tokenString,
//			401,
//			"Unauthorized",
//		},
//	}
//
//	// test each sample request payload
//	for _, v := range userSample {
//		// build the request
//		req, err := http.NewRequest("DELETE", "/users", nil)
//		if err != nil {
//			t.Errorf("Failed to create request: %v\n", err)
//		}
//		// set request variables and create response recorder
//		req = mux.SetURLVars(req, map[string]string{"id": v.id})
//		rr := httptest.NewRecorder()
//		handler := http.HandlerFunc(server.DeleteUserById)
//		// set token header
//		req.Header.Set("Authorization", v.tokenGiven)
//		// serve the request
//		handler.ServeHTTP(rr, req)
//
//		assert.Equal(t, rr.Code, v.stateCode)
//		// failed request tests
//		if v.stateCode == 401 && v.errorMessage != "" {
//			responseMap := make(map[string]interface{})
//			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
//			if err != nil {
//				t.Errorf("Cannot convert to json: %v\n", err)
//			}
//			assert.Equal(t, responseMap["error"], v.errorMessage)
//		}
//	}
//}
