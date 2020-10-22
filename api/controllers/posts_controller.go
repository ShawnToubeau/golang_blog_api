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

// Route handler for creating a post.
func (server *Server) InsertPost(w http.ResponseWriter, r *http.Request) {
	// read in request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	// extract post data from body
	post := models.Post{}
	err = json.Unmarshal(body, &post)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	// prepare and validate the post
	post.Prepare()
	err = post.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	// extract the user ID
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}
	// check if the user ID in the token matches the user ID in the post object
	if uid != post.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	// create the post
	postCreated, err := post.InsertPost(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	// set response header
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, postCreated.ID))
	responses.JSON(w, http.StatusOK, postCreated)
}

// Route handler for fetching all posts
func (server *Server) FetchAllPosts(w http.ResponseWriter, _ *http.Request) {
	// fetch posts
	post := models.Post{}
	posts, err := post.FetchAllPosts(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, posts)
}

// Route handler for fetching a post by a specific ID.
func (server *Server) GetPostByID(w http.ResponseWriter, r *http.Request) {
	// read in request variables
	vars := mux.Vars(r)
	// scan post ID
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	// fetch post
	post := models.Post{}
	postReceived, err := post.FetchPostById(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, postReceived)
}

// Route handler for updating a post by a specific ID.
func (server *Server) UpdatePostById(w http.ResponseWriter, r *http.Request) {
	// read in request variables
	vars := mux.Vars(r)
	// scan in post ID
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	// extract user ID
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	// fetch post
	post := models.Post{}
	err = server.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, err)
		return
	}
	// check if the request's user ID matches the user ID of the post
	if uid != post.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}
	// read in request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	// Assemble new post
	postUpdate := models.Post{}
	err = json.Unmarshal(body, &postUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	// check if the request's user ID matches the user ID of the new post
	if uid != postUpdate.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}
	// prepare and validate the new post
	postUpdate.Prepare()
	err = postUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	// update the post
	postUpdate.ID = post.ID
	postUpdated, err := postUpdate.UpdatePostById(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, postUpdated)
}

// Request handler for deleting a post by a specific ID.
func (server *Server) DeletePostById(w http.ResponseWriter, r *http.Request) {
	// read in request variables
	vars := mux.Vars(r)
	// scan post ID
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	// extract user ID
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	// fetch post
	post := models.Post{}
	err = server.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, err)
		return
	}
	// check if the request's user ID matches the user ID of the post
	if uid != post.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}
	// delete the post
	_, err = post.DeletePostById(server.DB, pid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	// set the header of the response
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")
}