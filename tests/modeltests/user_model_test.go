package modeltests

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/shawntoubeau/golang_blog_api/api/seed"
	"gopkg.in/go-playground/assert.v1"
	"testing"
)

// Fetch all users.
func TestFetchAllUsers(t *testing.T) {
	// seed test data
	users, _ := seed.Load(server.DB)
	// fetch users
	fetchedUsers, err := userInstance.FetchAllUsers(server.DB)
	if err != nil {
		t.Errorf("Failed fetching all users: %v\n", err)
		return
	}

	assert.Equal(t, len(*fetchedUsers), len(users))
}

// Insert a new user.
func TestInsertUser(t *testing.T) {
	// seed test data
	_, _ = seed.Load(server.DB)
	// create new user
	newUser := seed.GenerateNewUser("New User", "new@user.com", "password")
	savedUser, err := newUser.InsertUser(server.DB)
	if err != nil {
		t.Errorf("Failed creating user: %v", err)
		return
	}

	assert.Equal(t, newUser.ID, savedUser.ID)
	assert.Equal(t, newUser.Nickname, savedUser.Nickname)
	assert.Equal(t, newUser.Email, savedUser.Email)
}

// Fetch user by a specific ID.
func TestGetUserById(t *testing.T) {
	// seed test data
	users, _ := seed.Load(server.DB)
	// retrieve first user
	user := users[0]
	// fetch user
	foundUser, err := userInstance.FetchUserByID(server.DB, user.ID)
	if err != nil {
		t.Errorf("Failed to fetch user by ID: %v\n", err)
		return
	}
	assert.Equal(t, foundUser.ID, user.ID)
	assert.Equal(t, foundUser.Nickname, user.Nickname)
	assert.Equal(t, foundUser.Email, user.Email)
}

// Update user by specific ID.
func TestUpdateUserById(t *testing.T) {
	// seed test data
	users, _ := seed.Load(server.DB)
	// retrieve first user
	user := users[0]
	// update user
	user.Nickname = "Not Shawn"
	updatedUser, err := user.UpdateUserByID(server.DB, user.ID)
	if err != nil {
		t.Errorf("Error updating user by ID: %v", err)
		return
	}

	assert.Equal(t, updatedUser.ID, user.ID)
	assert.Equal(t, updatedUser.Email, user.Email)
	assert.Equal(t, updatedUser.Nickname, user.Nickname)
}

// Delete specific post by ID.
func TestDeleteUserById(t *testing.T) {
	// seed test data
	users, posts := seed.Load(server.DB)
	// delete all posts to avoid foreign key constraints
	for _, post := range posts {
		post.DeletePostById(server.DB, post.ID, post.AuthorID)
	}
	// retrieve first user
	user := users[0]
	// delete user
	isDeleted, err := userInstance.DeleteUserByID(server.DB, user.ID)
	if err != nil {
		t.Errorf("Error deleting user by ID: %v\n", err)
		return
	}

	assert.Equal(t, isDeleted, int64(1))
}
