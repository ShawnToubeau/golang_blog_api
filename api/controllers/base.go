package controllers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/shawntoubeau/golang_blog_api/api/models"
	"log"
	"net/http"
)

// Server object structure containing references to the database and router.
type Server struct {
	DB *gorm.DB
	Router *mux.Router
}

// Initialized a server instance using a database driver, user, password, port, hostname, and a database name.
func (server *Server) Initialize(DbDriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {
	var err error
	var DBURL string

	// MySQL database connection string
	if DbDriver == "mysql" {
		DBURL = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
	}

	// PostgreSQL database connection string
	if DbDriver == "postgres" {
		DBURL = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
	}

	// open a connection to the database and set the reference on the server object
	server.DB, err = gorm.Open(DbDriver, DBURL)
	if err != nil {
		fmt.Printf("Cannot connect to %s database", DbDriver)
		log.Fatalf("Connection error: %v\n", err)
	} else {
		fmt.Printf("We are connected to the %s database", DbDriver)
	}

	// migrate user and post model
	server.DB.Debug().AutoMigrate(&models.User{}, &models.Post{})
	// instantiate a new router
	server.Router = mux.NewRouter()
	// setup the routes
	server.initializeRoutes()
}

// Starts the server using the provided port.
func (server *Server) Run(addr string) {
	fmt.Printf("Listening on port %v\n", addr)
	// serve the routes on the provided port
	log.Fatal(http.ListenAndServe(addr, server.Router))
}