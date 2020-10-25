package modeltests

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gopkg.in/go-playground/assert.v1"
	"log"
	"testing"
)

// Fetch all users.
func TestFetchAllUsers(t *testing.T) {
	err := refreshTables()
	if err != nil {
		log.Fatalf("Failed to refresh tables: %v\n", err)
	}

	numUsers, err := seedUsers()
	if err != nil {
		log.Fatalf("Failed to seed post and user: %v\n", err)
	}

	users, err := userInstance.FetchAllUsers(server.DB)
	if err != nil {
		t.Errorf("Failed fetching all users: %v\n", err)
		return
	}

	assert.Equal(t, len(*users), len(numUsers))
}

// Insert a new user.
func TestInsertUser(t *testing.T) {
	err := refreshTables()
	if err != nil {
		log.Fatalf("Failed to refresh tables: %v\n", err)
	}

	newUser := MockUser1
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
	err := refreshTables()
	if err != nil {
		log.Fatalf("Failed to refresh tables: %v\n", err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Failing to seed user table: %v", err)
	}

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
	err := refreshTables()
	if err != nil {
		log.Fatalf("Failed to refresh tables: %v\n", err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Failed to see user: %v", err)
	}

	userUpdate := MockUser1
	userUpdate.Nickname = "Not Shawn"
	updatedUser, err := userUpdate.UpdateUserByID(server.DB, user.ID)
	if err != nil {
		t.Errorf("Error updating user by ID: %v", err)
		return
	}

	assert.Equal(t, updatedUser.ID, userUpdate.ID)
	assert.Equal(t, updatedUser.Email, userUpdate.Email)
	assert.Equal(t, updatedUser.Nickname, userUpdate.Nickname)
}

// Delete specific post by ID.
func TestDeleteUserById(t *testing.T) {
	err := refreshTables()
	if err != nil {
		log.Fatalf("Failed to refresh tables: %v\n", err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Error seeding user: %v\n", err)
	}

	isDeleted, err := userInstance.DeleteUserByID(server.DB, user.ID)
	if err != nil {
		t.Errorf("Error deleting user by ID: %v\n", err)
		return
	}

	assert.Equal(t, isDeleted, int64(1))
}
