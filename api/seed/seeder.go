package seed

import (
	"log"
	"github.com/jinzhu/gorm"
	"github.com/shawntoubeau/golang_blog_api/api/models"
)

// mock users
var users = []models.User{
	models.User{
		Nickname: "Shawn Toubeau",
		Email:    "shawn@aol.com",
		Password: "123",
	},
	models.User{
		Nickname: "Aria",
		Email:    "a@aol.com",
		Password: "321",
	},
}

// mock posts
var posts = []models.Post{
	models.Post{
		Title: "I like dogs",
		Content: "Dogs r cute",
	},
	models.Post{
		Title: "I like cats",
		Content: "Cats r cool",
	},
}

// Load in mock data.
func Load(db *gorm.DB) {
	// drop post and user tables if they exists
	err := db.Debug().DropTableIfExists(&models.Post{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	// create user and post tables
	err = db.Debug().AutoMigrate(&models.User{}, &models.Post{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}
	// add foreign key between author ID and user ID
	err = db.Debug().Model(&models.Post{}).AddForeignKey("author_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("cannot attach foreign key: %v", err)
	}

	// loop over mock users and add them to the database
	for i, _ := range users {
		// insert user
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		// set the author ID of corresponding post
		posts[i].AuthorID = users[i].ID
		// insert post with author ID link
		err = db.Debug().Model(&models.Post{}).Create(&posts[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}
	}
}