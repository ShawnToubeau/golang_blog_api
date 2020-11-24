package seed

import (
	"github.com/shawntoubeau/golang_blog_api/api/models"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"time"
)

// mock MockUsers
var MockUser1 = models.User{
	Nickname: "Swan",
	Email:    "shawn@aol.com",
	Password: "123",
}
var MockUser2 = models.User{
	Nickname: "Aria",
	Email:    "aria@aol.com",
	Password: "321",
}
var MockUsers = []models.User{
	MockUser1,
	MockUser2,
}

// mock MockPosts
var MockPost1 = models.Post{
	Title:   "I like dogs",
	Content: "Dogs are gr8 :)",
}
var MockPost2 = models.Post{
	Title:   "I like Cats",
	Content: "Cats are gr8 :)",
}
var MockPosts = []models.Post{
	MockPost1,
	MockPost2,
}

func GetPostsAuthorsPassword(authorId uint32) string {
	for _, user := range MockUsers {
		if user.ID == authorId {
			return user.Password
		}
	}

	return ""
}

func GenerateNewUser(nickname string, email string, password string) models.User {
	return models.User{
		Nickname: nickname,
		Email:    email,
		Password: password,
	}
}

func GenerateNewPost(title string, content string, userId uint32) models.Post {
	return models.Post{
		Title:    title,
		Content:  content,
		AuthorID: userId,
	}
}

// Load in mock data.
func Load(db *gorm.DB) ([]models.User, []models.Post) {
	var insertedUsers []models.User
	var insertedPosts []models.Post
	// drop tables if they exist
	db.Debug().Migrator().DropTable(&models.User{})
	db.Debug().Migrator().DropTable(&models.Post{})
	// migrate tables
	db.Debug().AutoMigrate(&models.User{}, &models.Post{})

	// insert MockUsers
	for _, user := range MockUsers {
		// encrypt user's password
		user.BeforeSave()
		// insert user
		err := db.Debug().Create(&user).Error
		if err != nil {
			log.Fatalf("cannot seed user: %v - user: %v\n", err, user)
		}

		insertedUsers = append(insertedUsers, user)
	}

	// get total # of MockUsers
	var numUsers int64
	db.Model(&MockUsers).Count(&numUsers)

	// insert MockPosts when we have MockUsers
	if numUsers > 0 {
		for _, post := range MockPosts {
			// fetch random user ID
			s1 := rand.NewSource(time.Now().UnixNano())
			r1 := rand.New(s1)
			var userId []uint32
			err := db.Select("id").Model(&models.User{}).Offset(r1.Intn(int(numUsers))).Limit(1).Take(&userId).Error
			if err != nil {
				log.Printf("failed to fetch user ID: %v\n", err)
			}

			// set author ID
			post.AuthorID = userId[0]
			// insert post with author ID link
			err = db.Debug().Create(&post).Error
			if err != nil {
				log.Fatalf("cannot seed post: %v - post %v\n", err, post)
			}

			insertedPosts = append(insertedPosts, post)
		}
	}

	return insertedUsers, insertedPosts
}
