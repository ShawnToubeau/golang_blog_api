package controllertests

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/shawntoubeau/golang_blog_api/api/controllers"
	"github.com/shawntoubeau/golang_blog_api/api/models"
	"log"
	"os"
	"testing"
)

var server = controllers.Server{}
var userInstance = models.User{}
var postInstance = models.Post{}

func Database() {
	var err error

	TestDbDriver := os.Getenv("TestDbDriver")

	if TestDbDriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("TestDbUser"), os.Getenv("TestDbPassword"), os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbName"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("Connected to the %s database\n", TestDbDriver)
		}
	}
	if TestDbDriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbUser"), os.Getenv("TestDbName"), os.Getenv("TestDbPassword"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("Connected to the %s database\n", TestDbDriver)
		}
	}
}

func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		log.Fatalf("Error loading env variables: %v\n", err)
	}
	Database()

	os.Exit(m.Run())
}

func refreshUserTable() error {
	err := server.DB.DropTableIfExists(&models.User{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refresed table")
	return nil
}

func seedOneUser() (models.User, error) {
	_ = refreshUserTable()

	user := models.User{
		Nickname: "Shawn",
		Email:    "shawn@aol.com",
		Password: "123",
	}

	err := server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		log.Fatalf("Cannot seed user table: %v", err)
	}
	return user, nil
}

func seedUsers() ([]models.User, error) {
	users := []models.User{
		models.User{
			Nickname: "Shawn",
			Email:    "shawn@aol.com",
			Password: "123",
		},
		models.User{
			Nickname: "Aria",
			Email:    "aria@aol.com",
			Password: "321",
		},
	}

	for i, _ := range users {
		err := server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			return []models.User{}, err
		}
	}

	return users, nil
}

func refreshUserAndPostTable() error {
	err := server.DB.DropTableIfExists(&models.User{}, &models.Post{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}, &models.Post{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed user and post tables")
	return nil
}

func seedOneUserAndOnePost() (models.Post, error) {
	err := refreshUserAndPostTable()
	if err != nil {
		return models.Post{}, err
	}
	user := models.User{
		Nickname: "Aria",
		Email:    "aria@aol.com",
		Password: "321",
	}
	err = server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.Post{}, err
	}
	post := models.Post{
		Title:    "I like dogs",
		Content:  "Dogs are gr8",
		AuthorID: user.ID,
	}
	err = server.DB.Model(&post).Create(&post).Error
	if err != nil {
		return models.Post{}, err
	}
	return post, nil
}

func seedUsersAndPosts() ([]models.User, []models.Post, error) {
	err := refreshUserAndPostTable()
	if err != nil {
		return []models.User{}, []models.Post{}, err
	}
	var users = []models.User{
		models.User{
			Nickname: "Shawn",
			Email:    "shawn@aol.com",
			Password: "123",
		},
		models.User{
			Nickname: "Aria",
			Email:    "aria@aol.com",
			Password: "321",
		},
	}

	var posts = []models.Post{
		models.Post{
			Title:    "I like cats",
			Content:  "Cats are gr8",
		},
		models.Post{
			Title:    "I like dogs",
			Content:  "Dogs are gr8",
		},
	}

	for i, _ := range users {
		err = server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("Cannot seed users table: %v", err)
		}
		posts[i].AuthorID = users[i].ID

		err = server.DB.Model(&models.Post{}).Create(&posts[i]).Error
		if err != nil {
			log.Fatalf("Cannot seed posts table: %v", err)
		}
	}
	return users, posts, nil
}