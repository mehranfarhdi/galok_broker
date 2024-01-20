package models

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // Import the SQLite driver for testing
)

// TestModels is a test suite for the models package.
func TestModels(t *testing.T) {
	// Connect to an SQLite in-memory database for testing
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Run migrations to create the 'users' table
	db.AutoMigrate(&User{})

	t.Run("TestUserCRUD", func(t *testing.T) {
		// Test user data
		testUser := User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "testpassword",
		}

		// Test Create
		createdUser, err := testUser.SaveUser(db)
		if err != nil {
			t.Fatalf("Error creating user: %v", err)
		}
		if createdUser.ID == 0 {
			t.Error("Expected user ID to be non-zero after creation")
		}

		// Test Read
		foundUser, err := FindUserByID(db, createdUser.ID)
		if err != nil {
			t.Fatalf("Error finding user: %v", err)
		}
		if foundUser.Username != testUser.Username || foundUser.Email != testUser.Email {
			t.Error("Expected retrieved user to match test user data")
		}

		// Test Update
		newUsername := "newtestuser"
		createdUser.Username = newUsername
		updatedUser, err := createdUser.UpdateAUser(db, createdUser.ID)
		if err != nil {
			t.Fatalf("Error updating user: %v", err)
		}
		if updatedUser.Username != newUsername {
			t.Error("Expected username to be updated")
		}

		// Test Delete
		rowsAffected, err := createdUser.DeleteAUser(db, createdUser.ID)
		if err != nil {
			t.Fatalf("Error deleting user: %v", err)
		}
		if rowsAffected != 1 {
			t.Error("Expected one row to be affected by delete operation")
		}

		// Test NotFound error
		notFoundUser, err := FindUserByID(db, createdUser.ID)
		if err == nil || !errors.Is(err, gorm.ErrRecordNotFound) {
			t.Error("Expected ErrRecordNotFound error")
		}
		if notFoundUser.ID != 0 {
			t.Error("Expected notFoundUser to have zero ID")
		}
	})

	t.Run("TestUserValidation", func(t *testing.T) {
		// Test user data with invalid email
		invalidEmailUser := User{
			Username: "testuser",
			Email:    "invalidemail",
			Password: "testpassword",
		}

		// Validate with "create" action should return an error
		err := invalidEmailUser.Validate("create")
		if err == nil {
			t.Error("Expected validation error for invalid email")
		}

		// Test user data with missing username
		missingUsernameUser := User{
			Email:    "test@example.com",
			Password: "testpassword",
		}

		// Validate with "create" action should return an error
		err = missingUsernameUser.Validate("create")
		if err == nil {
			t.Error("Expected validation error for missing username")
		}

		// Validate with "update" action should return an error for missing password
		err = missingUsernameUser.Validate("update")
		if err == nil {
			t.Error("Expected validation error for missing password during update")
		}
	})

	t.Run("TestPasswordHashing", func(t *testing.T) {
		// Test password hashing
		password := "testpassword"
		hashedPassword, err := Hash(password)
		if err != nil {
			t.Fatalf("Error hashing password: %v", err)
		}

		// Test VerifyPassword with correct password
		err = VerifyPassword(string(hashedPassword), password)
		if err != nil {
			t.Fatalf("Error verifying correct password: %v", err)
		}

		// Test VerifyPassword with incorrect password
		err = VerifyPassword(string(hashedPassword), "incorrectpassword")
		if err == nil {
			t.Error("Expected error for incorrect password")
		}
	})
}

// TestMain is a special test function that runs before other tests.
func TestMain(m *testing.M) {
	// Run tests
	exitCode := m.Run()

	// Exit with the test result
	fmt.Println("Test result code:", exitCode)
}
