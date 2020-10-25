package models

import (
	"errors"
	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"html"
	"log"
	"strings"
	"time"
)

// User object
type User struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Nickname  string    `gorm:"size:255;not null;unique" json:"nickname"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Hash the user's password.
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// Compares a hashed password against a password stored in the database.
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// Hashes a user's password right before they are inserted/updated in the database.
func (u *User) BeforeSave() error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// Prepares a user object before use. Trims white space from nickname and email.
// Sets created and updated times to the current time.
func (u *User) Prepare() {
	u.ID = 0
	u.Nickname = html.EscapeString(strings.TrimSpace(u.Nickname))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

// Check to see if nickname isn't empty.
func CheckNickname(nickname string) error {
	if nickname == "" {
		return errors.New("nickname required")
	}

	return nil
}

// Check to see if password isn't empty.
func CheckPassword(password string) error {
	if password == "" {
		return errors.New("password required")
	}

	return nil
}

// Check to see if email isn't empty and is a valid format.
func CheckEmail(email string) error {
	if email == "" {
		return errors.New("email required")
	}

	if err := checkmail.ValidateFormat(email); err != nil {
		return errors.New("invalid email")
	}

	return nil
}

// Helper function which validates user fields. Pre-defined validation sequences:
//
// 1. update - nickname, password, email
//
// 2. login - password, email
//
// 3. default - nickname, password, email
func (u *User) Validate(action string) error {
	var err error

	switch strings.ToLower(action) {
	case "update":
		if err = CheckNickname(u.Nickname); err != nil {
			return err
		}
		if err = CheckPassword(u.Password); err != nil {
			return err
		}
		if err = CheckEmail(u.Email); err != nil {
			return err
		}

		return nil

	case "login":
		if err = CheckPassword(u.Password); err != nil {
			return err
		}
		if err = CheckEmail(u.Email); err != nil {
			return err
		}

		return nil

	default:
		if err = CheckNickname(u.Nickname); err != nil {
			return err
		}
		if err = CheckPassword(u.Password); err != nil {
			return err
		}
		if err = CheckEmail(u.Email); err != nil {
			return err
		}

		return nil
	}
}

// Creates a new user entry.
func (u *User) InsertUser(db *gorm.DB) (*User, error) {
	err := db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}

	return u, nil
}

// Fetch all users. Limit to first 100.
func (u *User) FetchAllUsers(db *gorm.DB) (*[]User, error) {
	var users []User
	err := db.Debug().Model(&User{}).Limit(100).Find(&users).Error
	if err != nil {
		return &[]User{}, err
	}

	return &users, err
}

// Fetch user by a specific ID.
func (u *User) FetchUserByID(db *gorm.DB, uid uint32) (*User, error) {
	var err error
	err = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("user not found")
	}

	return u, err
}

// Update a user entry by a specific ID.
func (u *User) UpdateUserByID(db *gorm.DB, uid uint32) (*User, error) {
	// hash password
	err := u.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}
	// update user fields
	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumn(
		map[string]interface{}{
			"password":   u.Password,
			"nickname":   u.Nickname,
			"email":      u.Email,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &User{}, db.Error
	}
	// fetch newly updated user
	err = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}

	return u, nil
}

// Delete a user entry by a specific ID.
func (u *User) DeleteUserByID(db *gorm.DB, uid uint32) (int64, error) {
	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})
	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}
