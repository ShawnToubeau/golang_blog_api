package models

import (
	"errors"
	"gorm.io/gorm"
	"html"
	"strings"
	"time"
)

// Post object
type Post struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Title     string    `gorm:"size:255;not null;unique" json:"title"`
	Content   string    `gorm:"size:255;not null;" json:"content"`
	Author    User      `json:"author"`
	AuthorID  uint32    `gorm:"not null" json:"author_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Prepares the post object before use. Trims white space from title and content.
// Sets created and updated times to the current time.
func (p *Post) Prepare() {
	p.ID = 0
	p.Title = html.EscapeString(strings.TrimSpace(p.Title))
	p.Content = html.EscapeString(strings.TrimSpace(p.Content))
	p.Author = User{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

// Validates post objects for a title, content, and an author ID.
func (p *Post) Validate() error {
	if p.Title == "" {
		return errors.New("title required")
	}
	if p.Content == "" {
		return errors.New("content required")
	}
	if p.AuthorID < 1 {
		return errors.New("author ID required")
	}
	return nil
}

// Create a new post entry. Returns the created post and linked author.
func (p *Post) InsertPost(db *gorm.DB) (*Post, error) {
	// create post
	err := db.Debug().Model(&Post{}).Create(&p).Error
	if err != nil {
		return &Post{}, err
	}
	// fetch the user associated with the new post and set the author data on the post object
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Post{}, err
		}
	}
	return p, nil
}

// Fetches all posts. Limit to first 100.
func (p *Post) FetchAllPosts(db *gorm.DB) (*[]Post, error) {
	var posts []Post
	err := db.Debug().Model(&Post{}).Limit(100).Find(&posts).Error
	if err != nil {
		return &[]Post{}, err
	}
	if len(posts) > 0 {
		// loop over all posts and fetch their associated authors
		for i, _ := range posts {
			err := db.Debug().Model(&User{}).Where("id = ?", posts[i].AuthorID).Take(&posts[i].Author).Error
			if err != nil {
				return &[]Post{}, err
			}
		}
	}
	return &posts, nil
}

// Fetch post by a specific ID.
func (p *Post) FetchPostById(db *gorm.DB, pid uint64) (*Post, error) {
	err := db.Debug().Model(&Post{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Post{}, err
	}
	if p.ID != 0 {
		// fetch the post's author info
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Post{}, err
		}
	}
	return p, nil
}

// Update post by a specific ID.
func (p *Post) UpdatePostById(db *gorm.DB) (*Post, error) {
	// update the post's title and content
	err := db.Debug().Model(&Post{}).Where("id = ?", p.ID).Updates(
		Post{
			Title:     p.Title,
			Content:   p.Content,
			UpdatedAt: time.Now(),
		}).Error
	if err != nil {
		return &Post{}, err
	}

	if p.ID != 0 {
		// fetch the post's author info
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Post{}, err
		}
	}
	return p, nil
}

// Delete post by a specific ID.
func (p *Post) DeletePostById(db *gorm.DB, pid uint64, uid uint32) (int64, error) {
	// delete post with matching ID and author ID
	db = db.Debug().Model(&Post{}).Where("id = ? and author_id = ?", pid, uid).Take(&Post{}).Delete(&Post{})
	if db.Error != nil {
		if db.Error == gorm.ErrRecordNotFound {
			return 0, errors.New("post not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
