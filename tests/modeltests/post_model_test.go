package modeltests

import (
	"github.com/shawntoubeau/golang_blog_api/api/seed"
	"gopkg.in/go-playground/assert.v1"
	_ "gorm.io/driver/postgres"
	"testing"
)

// Fetch all posts.
func TestFetchAllPosts(t *testing.T) {
	// seed test data
	_, posts := seed.Load(server.DB)
	// fetch posts
	fetchedPosts, err := postInstance.FetchAllPosts(server.DB)
	if err != nil {
		t.Errorf("Failed to fetch posts: %v\n", err)
		return
	}
	assert.Equal(t, len(*fetchedPosts), len(posts))
}

// Insert a post.
func TestInsertPost(t *testing.T) {
	// seed test data
	users, _ := seed.Load(server.DB)
	// retrieve the first user
	user := users[0]
	// create and save post
	newPost := seed.GenerateNewPost("New post", "New Content", user.ID)
	savedPost, err := newPost.InsertPost(server.DB)
	if err != nil {
		t.Errorf("Failed to save post to database: %v\n", err)
		return
	}

	assert.Equal(t, newPost.ID, savedPost.ID)
	assert.Equal(t, newPost.Title, savedPost.Title)
	assert.Equal(t, newPost.Content, savedPost.Content)
	assert.Equal(t, newPost.AuthorID, savedPost.AuthorID)
}

// Fetch a specific post by it's ID.
func TestFetchPostById(t *testing.T) {
	// seed test data
	_, posts := seed.Load(server.DB)
	// retrieve first post
	post := posts[0]
	// fetch first post
	fetchedPost, err := postInstance.FetchPostById(server.DB, post.ID)
	if err != nil {
		t.Errorf("Failed to fetch post by ID: %v\n", err)
		return
	}

	assert.Equal(t, post.ID, fetchedPost.ID)
	assert.Equal(t, post.AuthorID, fetchedPost.AuthorID)
	assert.Equal(t, post.Title, fetchedPost.Title)
	assert.Equal(t, post.Content, fetchedPost.Content)
}

// Update a post by it's ID.
func TestUpdatePostById(t *testing.T) {
	// seed test data
	_, posts := seed.Load(server.DB)
	// retrieve first post
	post := posts[0]
	// edit post content
	post.Title = "New Title"
	post.Content = "New Content"
	// update post
	updatedPost, err := post.UpdatePostById(server.DB)
	if err != nil {
		t.Errorf("Failed to update post by ID: %v\n", err)
		return
	}
	assert.Equal(t, post.ID, updatedPost.ID)
	assert.Equal(t, post.AuthorID, updatedPost.AuthorID)
	assert.Equal(t, post.Content, updatedPost.Content)
	assert.Equal(t, post.Title, updatedPost.Title)
}

// Delete a specific post by it's ID.
func TestDeletePostById(t *testing.T) {
	// seed test data
	_, posts := seed.Load(server.DB)
	// retrieve first post
	post := posts[0]
	// delete post
	isDeleted, err := postInstance.DeletePostById(server.DB, post.ID, post.AuthorID)
	if err != nil {
		t.Errorf("Failed to delete user by ID: %v\n", err)
	}

	assert.Equal(t, isDeleted, int64(1))
}
