package controllertests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
	"github.com/shawntoubeau/golang_blog_api/api/models"
	_ "github.com/shawntoubeau/golang_blog_api/api/models"
	"github.com/shawntoubeau/golang_blog_api/api/seed"
	"gopkg.in/go-playground/assert.v1"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestCreatePost(t *testing.T) {
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

	// new post
	newPost := seed.GenerateNewPost("New Post", "New Post", user.ID)
	// sample request payloads
	validRequest := fmt.Sprintf(`{"title": "%v", "content": "%v", "author_id": %v}`, newPost.Title, newPost.Content, newPost.AuthorID)
	titleMissing := fmt.Sprintf(`{"title": "", "content": "%v", "author_id": %v}`, newPost.Content, newPost.AuthorID)
	contentMissing := fmt.Sprintf(`{"title": "%v", "content": "", "author_id": %v}`, newPost.Title, newPost.AuthorID)
	authorIDMissing := fmt.Sprintf(`{"title": "%v", "content": "%v", "author_id": 0}`, newPost.Title, newPost.Content)
	// sample request payloads and responses
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
			newPost.Title,
			newPost.Content,
			newPost.AuthorID,
			tokenString,
			"",
		},
		// Title missing
		{
			titleMissing,
			422,
			"",
			newPost.Content,
			newPost.AuthorID,
			tokenString,
			"title required",
		},
		// Content missing
		{
			contentMissing,
			422,
			newPost.Title,
			"",
			newPost.AuthorID,
			tokenString,
			"content required",
		},
		// Author ID missing
		{
			authorIDMissing,
			422,
			newPost.Title,
			newPost.Content,
			0,
			tokenString,
			"author ID required",
		},
		// Title taken
		{
			validRequest,
			500,
			newPost.Title,
			newPost.Content,
			newPost.AuthorID,
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
	// seed test data
	_, posts  := seed.Load(server.DB)
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
	var fetchedPosts []models.Post
	err = json.Unmarshal([]byte(rr.Body.String()), &fetchedPosts)
	if err != nil {
		log.Fatalf("Cannot convert to JSON: %v\n", err)
	}

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(fetchedPosts), len(posts))
}

func TestFetchPostById(t *testing.T) {
	// seed test data
	_, posts  := seed.Load(server.DB)
	// retrieve first post
	post := posts[0]
	// sample request payloads and responses
	samples := []struct {
		id         string
		statusCode int
		title      string
		content    string
		authorId   uint32
	}{
		{
			strconv.Itoa(int(post.ID)),
			200,
			post.Title,
			post.Content,
			post.AuthorID,
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
	// seed test data
	_, posts  := seed.Load(server.DB)
	// retrieve posts
	firstPost := posts[0]
	secondPost := posts[1]
	// retrieve first post's user
	user := firstPost.Author
	user.Password = seed.GetPostsAuthorsPassword(firstPost.AuthorID)

	// login user to retrieve auth token
	token, err := server.AuthenticateCredentials(user.Email, user.Password)
	if err != nil {
		log.Fatalf("Failed to login user: %v\n", err)
	}
	// build token string
	tokenString := fmt.Sprintf("Bearer %v", token)

	// sample request payloads
	validPayload := fmt.Sprintf(`{"title": "%v", "content": "%v", "author_id": %v}`, "new title", "new content", user.ID)
	missingTitle := fmt.Sprintf(`{"title": "", "content": "%v", "author_id": %v}`, "new content", user.ID)
	titleTaken := fmt.Sprintf(`{"title": "%v", "content": "%v", "author_id": %v}`, secondPost.Title, "new content", user.ID)
	missingContent := fmt.Sprintf(`{"title": "%v", "content": "", "author_id": %v}`, "new title", user.ID)

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
			strconv.Itoa(int(user.ID)),
			validPayload,
			200,
			"new title",
			"new content",
			tokenString,
			"",
		},
		// missing title
		{
			strconv.Itoa(int(user.ID)),
			missingTitle,
			422,
			"",
			"new content",
			tokenString,
			"title required",
		},
		// title taken
		{
			strconv.Itoa(int(user.ID)),
			titleTaken,
			500,
			posts[1].Title,
			"new content",
			tokenString,
			"title already taken",
		},
		// missing content
		{
			strconv.Itoa(int(user.ID)),
			missingContent,
			422,
			"new title",
			"",
			tokenString,
			"content required",
		},
		// missing token
		{
			strconv.Itoa(int(user.ID)),
			validPayload,
			401,
			"new title",
			"new content",
			"",
			"token contains an invalid number of segments",
		},
		// missing token
		{
			strconv.Itoa(int(user.ID)),
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
	// seed test data
	users, _  := seed.Load(server.DB)
	// retrieve first user
	user := users[0]
	user.Password = seed.MockUser1.Password
	// login the user to get their auth token
	token, err := server.AuthenticateCredentials(user.Email, user.Password)
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
			"token contains an invalid number of segments",
		},
		// Incorrect token string
		{
			strconv.Itoa(int(user.ID)),
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
			assert.Equal(t, responseMap["error"], sample.errorMessage)
		}
	}
}
