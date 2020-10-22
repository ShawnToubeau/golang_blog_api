package api

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/shawntoubeau/golang_blog_api/api/controllers"
	"github.com/shawntoubeau/golang_blog_api/api/seed"
	"log"
	"os"
)

var server = controllers.Server{}

// Start the server on a specified port.
func Run() {
	// load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env: %v", err)
	} else {
		fmt.Println("Pulling env variables")
	}
	// init server
	server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
	// seed database with mock data
	seed.Load(server.DB)
	// server on port
	server.Run(":8080")
}