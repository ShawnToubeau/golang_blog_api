package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/shawntoubeau/golang_blog_api/api/auth"
	"github.com/shawntoubeau/golang_blog_api/api/models"
	"github.com/shawntoubeau/golang_blog_api/api/responses"
	"github.com/shawntoubeau/golang_blog_api/api/utils/formaterror"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
)

// Sign in a user with an email and password. Returns an auth token using the user's ID.
func (server *Server) AuthenticateCredentials(email, password string) (string, error) {
	var err error
	// create user object
	user := models.User{}
	// fetch the user with the matching email provided
	err = server.DB.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return "", err
	}
	// verify that the password provided matches the password of the fetched user
	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	// create auth token
	return auth.CreateToken(user.ID)
}

// Route for processing login requests.
func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	// read request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	// create user object
	user := models.User{}
	// read user info from request body into user object
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user.Prepare()
	// todo: not sure what Prepare does here
	fmt.Printf("Login: user after prepare: %v\n", user)
	// validate the user for login
	err = user.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	// authenticate email and password and create an auth token
	token, err := server.AuthenticateCredentials(user.Email, user.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}
	// return auth token
	responses.JSON(w, http.StatusOK, token)
}
