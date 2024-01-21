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

// Package models provides data models and database operations for user management.
// It includes a User struct representing user data, functions for password hashing,
// validation, and CRUD operations on user records.
//Author : Mehran farhadi bajestani
//Created_at : 2024-01-20
//
// Structs:
//   - User: Represents a user with fields for ID, username, email, hashed password,
//           creation time, and last update time.
//
// Functions:
//   - Hash: Hashes a plain-text password using bcrypt.
//   - VerifyPassword: Compares a hashed password with a plain-text password.
//   - BeforeSave: Hashes the user's password before saving to the database.
//   - Prepare: Sanitizes and prepares user data before saving to the database.
//   - Validate: Validates user data based on the specified action (create, update, login, etc.).
//   - SaveUser: Creates a new user record in the database.
//   - FindAllUsers: Retrieves a list of all users from the database.
//   - FindUserByID: Retrieves a user by ID from the database.
//   - UpdateAUser: Updates a user's information in the database.
//   - DeleteAUser: Deletes a user from the database.
//
// Dependencies:
//   - github.com/badoux/checkmail: Used for email validation.
//   - github.com/jinzhu/gorm: ORM library for database interactions.
//   - golang.org/x/crypto/bcrypt: Used for password hashing.
//
// Note: Before using the SaveUser, UpdateAUser, and DeleteAUser functions,
// the database connection (gorm.DB) should be provided to these functions.

type User struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Username  string    `gorm:"size:255;not null;unique" json:"username"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	IsAdmin   bool      `gorm:"default:false;"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *User) BeforeSave() error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) Prepare() {
	u.ID = 0
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if u.Username == "" {
			return errors.New("Reuired Nicname")
		}
		if u.Password == "" {
			return errors.New("Password Required")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}

		return nil
	case "login_by_username":
		if u.Username == "" {
			return errors.New("Reuired Nicname")
		}
		if u.Password == "" {
			return errors.New("Password Required")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}

		return nil
	case "login_by_email":
		if u.Password == "" {
			return errors.New("Password Required")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}

		return nil
	default:
		if u.Username == "" {
			return errors.New("Reuired Nicname")
		}
		if u.Password == "" {
			return errors.New("Password Required")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}

		return nil
	}
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {

	var err error
	err = db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) FindAllUsers(db *gorm.DB) (*[]User, error) {
	var err error
	users := []User{}
	err = db.Debug().Model(&User{}).Limit(100).Find(&users).Error
	if err != nil {
		return &[]User{}, err
	}
	return &users, err
}

func (u *User) FindUserByID(db *gorm.DB, uid uint32) (*User, error) {
	var err error
	err = db.Debug().Model(User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("User Not Found")
	}
	return u, err
}

func (u *User) UpdateAUser(db *gorm.DB, uid uint32) (*User, error) {

	// To hash the password
	err := u.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}
	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"password":   u.Password,
			"username":   u.Username,
			"email":      u.Email,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &User{}, db.Error
	}
	// This is the display the updated user
	err = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) DeleteAUser(db *gorm.DB, uid uint32) (int64, error) {

	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
