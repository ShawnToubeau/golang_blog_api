package middlewares

import (
	"errors"
	"net/http"
	"github.com/shawntoubeau/golang_blog_api/api/auth"
	"github.com/shawntoubeau/golang_blog_api/api/responses"
)

// Sets the content type of all response writers to application/json.
func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

// Authentication middleware which validates the auth token provided in the request header.
func SetMiddlewareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.TokenValid(r)
		if err != nil {
			responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
			return
		}
		next(w, r)
	}
}