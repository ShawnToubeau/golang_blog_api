package modeltests

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/shawntoubeau/golang_blog_api/api/controllers"
	"github.com/shawntoubeau/golang_blog_api/api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"testing"
)

var server = controllers.Server{}
var userInstance = models.User{}
var postInstance = models.Post{}

// Open a connection to the test database
func Database() {
	var host, port, user, password, dbname string
	dbDriver := os.Getenv("TEST_DB_DRIVER")
	var err error

	if dbDriver == "postgres" {
		// set env variables
		host = os.Getenv("TEST_DB_HOST")
		port = os.Getenv("TEST_DB_PORT")
		user = os.Getenv("TEST_DB_USER")
		password = os.Getenv("TEST_DB_PASSWORD")
		dbname = os.Getenv("TEST_DB_NAME")

		DBURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable ", host, port, user, password, dbname)
		server.DB, err = gorm.Open(postgres.Open(DBURL), &gorm.Config{})
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", dbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("Connected to the %s database\n", dbDriver)
		}
	}
}

// Test entry point.
func TestMain(m *testing.M) {
	// load environment vars
	err := godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		log.Fatalf("Error loading env variables: %v\n", err)
	}
	// create database connection
	Database()
	// run the tests
	os.Exit(m.Run())
}
