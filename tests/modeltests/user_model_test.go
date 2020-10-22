package modeltests

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/shawntoubeau/golang_blog_api/api/models"
	"gopkg.in/go-playground/assert.v1"
	"log"
	"testing"
)

func TestFindAllUsers(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	numUsers, err := seedUsers()
	if err != nil {
		log.Fatal(err)
	}
	users, err := userInstance.FetchAllUsers(server.DB)
	if err != nil {
		t.Errorf("Failing fetching all users: %v\n", err)
		return
	}
	assert.Equal(t, len(*users), numUsers)
}

func TestSaveUser(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	newUser := models.User{
		ID:       1,
		Nickname: "test",
		Email:    "test@aol.com",
		Password: "123",
	}
	savedUser, err := newUser.InsertUser(server.DB)
	if err != nil {
		t.Errorf("Failed creating user: %v", err)
		return
	}
	assert.Equal(t, newUser.ID, savedUser.ID)
	assert.Equal(t, newUser.Nickname, savedUser.Nickname)
	assert.Equal(t, newUser.Email, savedUser.Email)
}

func TestGetUserById(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
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

func TestUpdateUserById(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Failed to see user: %v", err)
	}

	userUpdate := models.User{
		ID:       user.ID,
		Nickname: "New Test",
		Email:    "shawn@aol.com",
		Password: "123",
	}
	updatedUser, err := userUpdate.UpdateUserByID(server.DB, user.ID)
	if err != nil {
		t.Errorf("Error updating user by ID: %v", err)
		return
	}
	assert.Equal(t, updatedUser.ID, userUpdate.ID)
	assert.Equal(t, updatedUser.Email, userUpdate.Email)
	assert.Equal(t, updatedUser.Nickname, userUpdate.Nickname)
}

func TestDeleteUserById(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
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