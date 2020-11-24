package controllertests

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
	var err error

	TestDbDriver := os.Getenv("TestDbDriver")

	if TestDbDriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbUser"), os.Getenv("TestDbName"), os.Getenv("TestDbPassword"))
		server.DB, err = gorm.Open(postgres.Open(DBURL), &gorm.Config{})
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("Connected to the %s database\n", TestDbDriver)
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

// Drops old user and post tables and migrates user and post schemas.
//func refreshTables() error {
//	// drop tables if they exist
//	server.DB.Migrator().DropTable(&models.User{})
//	server.DB.Migrator().DropTable(&models.Post{})
//	// migrate tables
//	server.DB.AutoMigrate(&models.User{}, &models.Post{})
//
//	log.Printf("Successfully refreshed user and post tables")
//	return nil
//}
//
//// Insert 1 mock user into the database.
//func seedOneUser() (models.User, error) {
//	user := MockUser1
//	user.BeforeSave()
//	err := server.DB.Create(&user).Error
//	if err != nil {
//		return models.User{}, err
//	}
//
//	log.Printf("Successfully seeded user tables")
//	return user, nil
//}
//
//// Insert multiple mock users into the database.
//func seedUsers() ([]models.User, error) {
//	users := []models.User{
//		MockUser1,
//		MockUser2,
//	}
//
//	// insert users
//	for i, _ := range users {
//		err := server.DB.Model(&models.User{}).Create(&users[i]).Error
//		if err != nil {
//			return []models.User{}, err
//		}
//	}
//
//	log.Printf("Successfully seeded user tables")
//	return users, nil
//}
//
//// Insert 1 mock user and 1 mock post.
//func seedOneUserAndOnePost() (models.User, models.Post, error) {
//	user := MockUser1
//	err := server.DB.Model(&models.User{}).Create(&user).Error
//	if err != nil {
//		return models.User{}, models.Post{}, err
//	}
//
//	post := MockPost1(user.ID)
//	err = server.DB.Model(&post).Create(&post).Error
//	if err != nil {
//		return models.User{}, models.Post{}, err
//	}
//
//	log.Printf("Successfully seeded user and post tables")
//	return user, post, nil
//}
//
//// Insert multiple mock users and mock posts.
//func seedUsersAndPosts() ([]models.User, []models.Post, error) {
//	user1 := MockUser1
//	user2 := MockUser2
//	var users = []models.User{
//		user1,
//		user2,
//	}
//
//	// insert users
//	for i, _ := range users {
//		err := server.DB.Model(&models.User{}).Create(&users[i]).Error
//		if err != nil {
//			return []models.User{}, []models.Post{}, err
//		}
//	}
//
//	var posts = []models.Post{
//		MockPost1(users[0].ID),
//		MockPost2(users[1].ID),
//	}
//
//	// insert posts
//	for i, _ := range posts {
//		err := server.DB.Model(&models.Post{}).Create(&posts[i]).Error
//		if err != nil {
//			return []models.User{}, []models.Post{}, err
//		}
//	}
//
//	log.Printf("Successfully seeded user and post tables")
//	return users, posts, nil
//}
