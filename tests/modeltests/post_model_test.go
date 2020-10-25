package modeltests

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gopkg.in/go-playground/assert.v1"
	"log"
	"testing"
)

// Fetch all posts.
func TestFetchAllPosts(t *testing.T) {
	_, posts, err := seedUsersAndPosts()
	if err != nil {
		log.Fatalf("Failed to seed post and user: %v\n", err)
	}

	fetchedPosts, err := postInstance.FetchAllPosts(server.DB)
	if err != nil {
		t.Errorf("Failed to fetch posts: %v\n", err)
		return
	}
	assert.Equal(t, len(*fetchedPosts), len(posts))
}

// Insert a post.
func TestInsertPost(t *testing.T) {
	err := refreshTables()
	if err != nil {
		log.Fatalf("Failed to refresh tables: %v\n", err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Failed to save user to database: %v\n", err)
	}

	// create and save post
	newPost := MockPost1(user.ID)
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
	err := refreshTables()
	if err != nil {
		log.Fatalf("Failed to refresh tables: %v\n", err)
	}

	_, post, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatalf("Failed to seed user and post tables: %v\n", err)
	}

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
	err := refreshTables()
	if err != nil {
		log.Fatalf("Failed to refresh tables: %v\n", err)
	}

	_, post, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatalf("Failed to see user and post: %v\n", err)
	}

	editedPost := post
	editedPost.Title = "New Title"
	editedPost.Content = "New Content"

	updatedPost, err := editedPost.UpdatePostById(server.DB)
	if err != nil {
		t.Errorf("Failed to update post by ID: %v\n", err)
		return
	}
	assert.Equal(t, editedPost.ID, updatedPost.ID)
	assert.Equal(t, editedPost.AuthorID, updatedPost.AuthorID)
	assert.Equal(t, editedPost.Content, updatedPost.Content)
	assert.Equal(t, editedPost.Title, updatedPost.Title)
}

// Delete a specific post by it's ID.
func TestDeletePostById(t *testing.T) {
	err := refreshTables()
	if err != nil {
		log.Fatalf("Failed to refresh tables: %v\n", err)
	}

	_, post, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatalf("Failed to seed user and post: %v\n", err)
	}

	isDeleted, err := postInstance.DeletePostById(server.DB, post.ID, post.AuthorID)
	if err != nil {
		t.Errorf("Failed to delete user by ID: %v\n", err)
	}

	assert.Equal(t, isDeleted, int64(1))
}
