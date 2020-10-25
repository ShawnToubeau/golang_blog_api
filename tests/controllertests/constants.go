package controllertests

import (
	"github.com/shawntoubeau/golang_blog_api/api/models"
)

var MockUser1 = models.User{
	Nickname: "Shawn",
	Email:    "shawn@aol.com",
	Password: "123",
}

var MockUser2 = models.User{
	Nickname: "Aria",
	Email:    "aria@aol.com",
	Password: "321",
}

func MockPost1(userId uint32) models.Post {
	return models.Post{
		Title:    "I like dogs",
		Content:  "Dogs are gr8 :)",
		AuthorID: userId,
	}
}

func MockPost2(userId uint32) models.Post {
	return models.Post{
		Title:    "I like Cats",
		Content:  "Cats are gr8 :)",
		AuthorID: userId,
	}
}
