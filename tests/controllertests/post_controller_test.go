package controllertests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
	"github.com/shawntoubeau/golang_blog_api/api/models"
	_ "github.com/shawntoubeau/golang_blog_api/api/models"
	"gopkg.in/go-playground/assert.v1"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestCreatePost(t *testing.T) {
	var AuthEmail, AuthPassword string

	// refresh tables
	err := refreshTables()
	if err != nil {
		log.Fatalf("Failed to refresh tables: %v\n", err)
	}
	// seed users
	users, err := seedUsers()
	if err != nil {
		log.Fatalf("Failed to seed users")
	}
	// retrieve first user
	AuthEmail = users[0].Email
	AuthPassword = MockUser1.Password
	// login user to retrieve auth token
	token, err := server.AuthenticateCredentials(AuthEmail, AuthPassword)
	if err != nil {
		log.Fatalf("Failed to login user: %v\n", err)
	}
	// format token
	tokenString := fmt.Sprintf("Bearer %v", token)

	// mock posts
	post1 := MockPost1(users[0].ID)
	// mock request payloads
	validRequest := fmt.Sprintf(`{"title": "%v", "content": "%v", "author_id": %v}`, post1.Title, post1.Content, post1.AuthorID)
	titleMissing := fmt.Sprintf(`{"title": "", "content": "%v", "author_id": %v}`, post1.Content, post1.AuthorID)
	contentMissing := fmt.Sprintf(`{"title": "%v", "content": "", "author_id": %v}`, post1.Title, post1.AuthorID)
	authorIDMissing := fmt.Sprintf(`{"title": "%v", "content": "%v", "author_id": 0}`, post1.Title, post1.Content)

	samples := []struct {
		inputJSON    string
		statusCode   int
		title        string
		content      string
		authorID     uint32
		tokenGiven   string
		errorMessage string
	}{
		// Valid
		{
			validRequest,
			200,
			post1.Title,
			post1.Content,
			post1.AuthorID,
			tokenString,
			"",
		},
		// Title missing
		{
			titleMissing,
			422,
			"",
			post1.Content,
			post1.AuthorID,
			tokenString,
			"title required",
		},
		// Content missing
		{
			contentMissing,
			422,
			post1.Title,
			"",
			post1.AuthorID,
			tokenString,
			"content required",
		},
		// Author ID missing
		{
			authorIDMissing,
			422,
			post1.Title,
			post1.Content,
			0,
			tokenString,
			"author ID required",
		},
		// Title taken
		{
			validRequest,
			500,
			post1.Title,
			post1.Content,
			post1.AuthorID,
			tokenString,
			"title already taken",
		},
	}

	// test sample requests
	for _, sample := range samples {
		// build the request
		req, err := http.NewRequest("PUT", "/post", bytes.NewBufferString(sample.inputJSON))
		if err != nil {
			t.Errorf("Error creating request: %v\n", err)
		}
		// create response recorder and serve the request
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.InsertPost)
		req.Header.Set("Authorization", sample.tokenGiven)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to JSON: %v\n", err)
		}

		assert.Equal(t, rr.Code, sample.statusCode)
		// valid request tests
		if sample.statusCode == 200 {
			assert.Equal(t, responseMap["title"], sample.title)
			assert.Equal(t, responseMap["content"], sample.content)
			assert.Equal(t, responseMap["author_id"], float64(sample.authorID))
		}
		// invalid request tests
		if sample.statusCode == 422 || sample.statusCode == 500 {
			assert.Equal(t, responseMap["error"], sample.errorMessage)
		}
	}
}

func TestFetchPosts(t *testing.T) {
	// refresh tables
	err := refreshTables()
	if err != nil {
		log.Fatalf("Failed to refresh tables: %v\n", err)
	}
	// seed tables
	_, mockPosts, err := seedUsersAndPosts()
	if err != nil {
		log.Fatalf("Failed to seed tables: %v\n", err)
	}
	// create request
	req, err := http.NewRequest("GET", "/posts", nil)
	if err != nil {
		t.Errorf("Failed to create request: %v\n", err)
	}
	// create response recorder
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.FetchAllPosts)
	// serve request
	handler.ServeHTTP(rr, req)
	// create post array and process response
	var posts []models.Post
	err = json.Unmarshal([]byte(rr.Body.String()), &posts)
	if err != nil {
		log.Fatalf("Cannot convert to JSON: %v\n", err)
	}

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(posts), len(mockPosts))
}

func TestFetchPostById(t *testing.T) {
	// refresh tables
	err := refreshTables()
	if err != nil {
		log.Fatalf("Failed to refresh tables: %v\n", err)
	}
	// seed tables
	_, posts, err := seedUsersAndPosts()
	if err != nil {
		log.Fatalf("Failed to seed tables: %v\n", err)
	}
	// sample request payloads and responses
	samples := []struct {
		id         string
		statusCode int
		title      string
		content    string
		authorId   uint32
	}{
		{
			strconv.Itoa(int(posts[0].ID)),
			200,
			posts[0].Title,
			posts[0].Content,
			posts[0].AuthorID,
		},
		{
			"unknown",
			400,
			"",
			"",
			0,
		},
	}

	for _, sample := range samples {
		// create request
		req, err := http.NewRequest("GET", "/posts", nil)
		if err != nil {
			t.Errorf("Failed to create request: %v\n", err)
		}
		// set request params
		req = mux.SetURLVars(req, map[string]string{"id": sample.id})
		// create response recorder
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.FetchPostByID)
		// serve the request
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
			assert.Equal(t, responseMap["title"], sample.title)
			assert.Equal(t, responseMap["content"], sample.content)
			assert.Equal(t, responseMap["author_id"], float64(sample.authorId))
		}
	}
}

func TestUpdatePostById(t *testing.T) {
	var AuthEmail, AuthPassword string
	var AuthID uint32
	// refresh tables
	err := refreshTables()
	if err != nil {
		log.Fatalf("Failed to refresh tables: %v\n", err)
	}
	// seed tables
	users, posts, err := seedUsersAndPosts()
	if err != nil {
		log.Fatalf("Failed to seed tables: %v\n", err)
	}
	fmt.Printf("Users: %v\n", users)
	fmt.Printf("Posts: %v\n", posts)
	// retrieve first post's user
	AuthID = posts[0].AuthorID
	AuthEmail = users[0].Email
	AuthPassword = MockUser1.Password
	fmt.Printf("User creds: %v %v\n", AuthEmail, AuthPassword)

	// login user to retrieve auth token
	token, err := server.AuthenticateCredentials(AuthEmail, AuthPassword)
	if err != nil {
		log.Fatalf("Failed to login user: %v\n", err)
	}
	// build token string
	tokenString := fmt.Sprintf("Bearer %v", token)

	// mock request JSON payloads
	validPayload := fmt.Sprintf(`{"title": "%v", "content": "%v", "author_id": %v}`, "new title", "new content", AuthID)
	missingTitle := fmt.Sprintf(`{"title": "", "content": "%v", "author_id": %v}`, "new content", AuthID)
	titleTaken := fmt.Sprintf(`{"title": "%v", "content": "%v", "author_id": %v}`, posts[1].Title, "new content", AuthID)
	missingContent := fmt.Sprintf(`{"title": "%v", "content": "", "author_id": %v}`, "new title", AuthID)

	// sample request payloads and responses
	samples := []struct {
		id            string
		updateJSON    string
		statusCode    int
		updateTitle   string
		updateContent string
		tokenGiven    string
		errorMessage  string
	}{
		// valid case
		{
			strconv.Itoa(int(AuthID)),
			validPayload,
			200,
			"new title",
			"new content",
			tokenString,
			"",
		},
		// missing title
		{
			strconv.Itoa(int(AuthID)),
			missingTitle,
			422,
			"",
			"new content",
			tokenString,
			"title required",
		},
		// title taken
		{
			strconv.Itoa(int(AuthID)),
			titleTaken,
			500,
			posts[1].Title,
			"new content",
			tokenString,
			"title already taken",
		},
		// missing content
		{
			strconv.Itoa(int(AuthID)),
			missingContent,
			422,
			"new title",
			"",
			tokenString,
			"content required",
		},
		// missing token
		{
			strconv.Itoa(int(AuthID)),
			validPayload,
			401,
			"new title",
			"new content",
			"",
			"token contains an invalid number of segments",
		},
		// missing token
		{
			strconv.Itoa(int(AuthID)),
			validPayload,
			401,
			"new title",
			"new content",
			"incorrect token",
			"token contains an invalid number of segments",
		},
		// missing author ID
		{
			"",
			validPayload,
			400,
			"new title",
			"new content",
			tokenString,
			"strconv.ParseUint: parsing \"\": invalid syntax",
		},
		// incorrect author ID
		{
			strconv.Itoa(2),
			validPayload,
			401,
			"new title",
			"new content",
			tokenString,
			"unauthorized",
		},
	}

	// test sample requests
	for _, sample := range samples {
		// create request
		req, err := http.NewRequest("POST", "/posts", bytes.NewBufferString(sample.updateJSON))
		if err != nil {
			t.Errorf("Error creating request: %v\n", err)
		}
		// set request params
		req = mux.SetURLVars(req, map[string]string{"id": sample.id})
		req.Header.Set("Authorization", sample.tokenGiven)
		// create response recorder
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdatePostById)
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
			assert.Equal(t, responseMap["title"], sample.updateTitle)
			assert.Equal(t, responseMap["content"], sample.updateContent)
		}
		// invalid request tests
		if sample.statusCode != 200 {
			assert.Equal(t, responseMap["error"], sample.errorMessage)
		}
	}
}

func TestDeletePost(t *testing.T) {
	var AuthEmail, AuthPassword string
	var AuthID uint32
	// refresh tables
	err := refreshTables()
	if err != nil {
		log.Fatalf("Failed to refresh tables: %v\n", err)
	}
	// seed tables
	users, _, err := seedUsersAndPosts()
	if err != nil {
		log.Fatalf("Failed to seed tables: %v\n", err)
	}
	// get first users credentials
	AuthID = users[0].ID
	AuthEmail = users[0].Email
	AuthPassword = MockUser1.Password
	// login the user to get their auth token
	token, err := server.AuthenticateCredentials(AuthEmail, AuthPassword)
	if err != nil {
		log.Fatalf("Failed to login user: %v\n", err)
	}
	// construct token string
	tokenString := fmt.Sprintf("Bearer: %v", token)

	// sample request payloads and responses
	userSamples := []struct {
		id           string
		tokenGiven   string
		statusCode   int
		errorMessage string
	}{
		// Valid
		{
			strconv.Itoa(int(AuthID)),
			tokenString,
			204,
			"",
		},
		// Missing token string
		{
			strconv.Itoa(int(AuthID)),
			"",
			401,
			"token contains an invalid number of segments",
		},
		// Incorrect token string
		{
			strconv.Itoa(int(AuthID)),
			"incorrect token string",
			401,
			"token contains an invalid number of segments",
		},
		// Missing post ID
		{
			"",
			tokenString,
			400,
			"strconv.ParseUint: parsing \"\": invalid syntax",
		},
		// Wrong post ID
		{
			strconv.Itoa(2),
			tokenString,
			401,
			"unauthorized",
		},
	}

	// test each sample request
	for _, sample := range userSamples {
		// build the request
		req, err := http.NewRequest("DELETE", "/posts", nil)
		if err != nil {
			t.Errorf("Failed to create request: %v\n", err)
		}
		// set request variables and create response recorder
		req = mux.SetURLVars(req, map[string]string{"id": sample.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeletePostById)
		// set token header
		req.Header.Set("Authorization", sample.tokenGiven)
		// serve the request
		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, sample.statusCode)
		if sample.statusCode != 204 && sample.errorMessage != "" {
			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				t.Errorf("Cannot convert to json: %v\n", err)
			}
			fmt.Printf("res map: %v\n", responseMap)
			assert.Equal(t, responseMap["error"], sample.errorMessage)
		}
	}
}
