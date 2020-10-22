package modeltests

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/shawntoubeau/golang_blog_api/api/models"
	"gopkg.in/go-playground/assert.v1"
	"log"
	"testing"
)

func TestFindAllPosts(t *testing.T) {
	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Failed refreshing post and user tables: %v\n", err)
	}
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

func TestSavePost(t *testing.T) {
	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Failed refreshing post and user tables: %v\n", err)
	}
	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Failed to save user to database: %v\n", err)
	}
	newPost := models.Post{
		ID:       1,
		Title:    "Test Title",
		Content:  "Test Content",
		AuthorID: user.ID,
	}
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

func TestGetPostById(t *testing.T) {
	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Failed refreshing post and user tables: %v\n", err)
	}
	post, err := seedOneUserAndOnePost()
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

func TestUpdatePostById(t *testing.T) {
	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Failed refreshing post and user tables: %v\n", err)
	}
	post, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatalf("Failed to see user and post: %v\n", err)
	}
	postUpdate := models.Post{
		ID:       post.ID,
		Title:    "New Title",
		Content:  "New Content",
		AuthorID: post.AuthorID,
	}
	updatedPost, err := postUpdate.UpdatePostById(server.DB)
	if err != nil {
		t.Errorf("Failed to update post by ID: %v\n", err)
		return
	}
	assert.Equal(t, postUpdate.ID, updatedPost.ID)
	assert.Equal(t, postUpdate.AuthorID, updatedPost.AuthorID)
	assert.Equal(t, postUpdate.Content, updatedPost.Content)
	assert.Equal(t, postUpdate.Title, updatedPost.Title)
}

func TestDeletePostById(t *testing.T) {
	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Failed refreshing post and user table: %v\n", err)
	}
	post, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatalf("Failed to seed user and post: %v\n", err)
	}
	isDeleted, err := postInstance.DeletePostById(server.DB, post.ID, post.AuthorID)
	if err != nil {
		t.Errorf("Failed to delete user by ID: %v\n", err)
	}

	assert.Equal(t, isDeleted, int64(1))
}