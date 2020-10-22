package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/shawntoubeau/golang_blog_api/api/auth"
	"github.com/shawntoubeau/golang_blog_api/api/models"
	"github.com/shawntoubeau/golang_blog_api/api/responses"
	"github.com/shawntoubeau/golang_blog_api/api/utils/formaterror"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Route handler for creating a user.
func (server *Server) InsertUser(w http.ResponseWriter, r *http.Request) {
	// read in request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	// extract user data from body
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	// prepare and validate the user
	user.Prepare()
	err = user.Validate("")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	// insert the user
	userCreated, err := user.InsertUser(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	// set response header
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, userCreated.ID))
	responses.JSON(w, http.StatusCreated, userCreated)
}

// Route handler for fetching all users.
func (server *Server) FetchAllUsers(w http.ResponseWriter, r *http.Request) {
	// create user model interface
	user := models.User{}
	// fetch all users
	users, err := user.FetchAllUsers(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, users)
}

// Route handler for fetching a specific user by ID.
func (server *Server) GetUserByID(w http.ResponseWriter, r *http.Request) {
	// read in request variables
	vars := mux.Vars(r)
	// parse the user ID
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	// create user model interface
	user := models.User{}
	// fetch the user using the user ID
	fetchedUser, err := user.FetchUserByID(server.DB, uint32(uid))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, fetchedUser)
}

// Route handler for updating a specific user by ID.
func (server *Server) UpdateUserById(w http.ResponseWriter, r *http.Request) {
	// read in request variables
	vars := mux.Vars(r)
	// parse the user ID
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	// read request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	// create user model interface and read user data into object
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	// extract auth token
	tokenId, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	// check if auth token matches the user token in the request
	if tokenId != uint32(uid) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	// prepare and validate the user
	user.Prepare()
	err = user.Validate("update")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	// update the user
	updatedUser, err := user.UpdateUserByID(server.DB, uint32(uid))
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, updatedUser)
}

// Route handler for deleting a specific user by ID.
func (server *Server) DeleteUserById(w http.ResponseWriter, r *http.Request) {
	// read in request variables
	vars := mux.Vars(r)
	// create user model interface
	user := models.User{}
	// parse user ID
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	// extract auth token ID
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorizerd"))
		return
	}
	// verify the auth token ID matches the user ID
	if tokenID != 0 && tokenID != uint32(uid) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	// delete the user
	_, err = user.DeleteUserByID(server.DB, uint32(uid))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	// set response header
	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
	responses.JSON(w, http.StatusNoContent, "")
}
