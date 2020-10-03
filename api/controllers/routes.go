package controllers

import "github.com/shawntoubeau/golang_blog_api/api/middlewares"

func (server *Server) initializeRoutes() {
	// Home
	server.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(server.Home)).Methods("GET")
	// Login
	server.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(server.Login)).Methods("POST")
	// User
	server.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(server.CreateUser)).Methods("POST")
	server.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(server.GetUsers)).Methods("GET")
	server.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(server.GetUserByID)).Methods("GET")
	server.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(server.UpdateUserById))).Methods("PUT")
	server.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(server.DeleteUserById)).Methods("DELETE")
	// Posts
	server.Router.HandleFunc("/posts", middlewares.SetMiddlewareJSON(server.CreatePost)).Methods("POST")
	server.Router.HandleFunc("/posts", middlewares.SetMiddlewareJSON(server.GetPosts)).Methods("GET")
	server.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareJSON(server.GetPostByID)).Methods("GET")
	server.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(server.UpdatePostById))).Methods("PUT")
	server.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareAuthentication(server.DeletePostById)).Methods("DELETE")
}


