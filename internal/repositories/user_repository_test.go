package repositories_test

import (
	"context"
	"database/sql"
	"mordezzanV4/internal/repositories"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) (*sql.DB, func()) {
	// Create a temporary database for testing
	dbFile := "./test_user_repo.db"
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Create the necessary table
	_, err = db.Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	// Return the DB and a cleanup function
	return db, func() {
		db.Close()
		os.Remove(dbFile)
	}
}

func TestUserRepositoryCRUD(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Initialize the repository
	repo := repositories.NewSQLCUserRepository(db)

	// Test data
	ctx := context.Background()
	testUsername := "testuser"
	testEmail := "test@example.com"
	testPassword := "hashedpassword123"

	t.Run("Create user", func(t *testing.T) {
		id, err := repo.CreateUser(ctx, testUsername, testEmail, testPassword)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}
		if id <= 0 {
			t.Fatalf("Expected valid user ID, got %d", id)
		}
	})

	t.Run("Get user", func(t *testing.T) {
		// First create a user to retrieve
		id, err := repo.CreateUser(ctx, "getuser", "get@example.com", testPassword)
		if err != nil {
			t.Fatalf("Failed to create user for get test: %v", err)
		}

		// Retrieve the user
		user, err := repo.GetUser(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get user: %v", err)
		}

		// Validate user data
		if user.ID != id {
			t.Errorf("Expected user ID %d, got %d", id, user.ID)
		}
		if user.Username != "getuser" {
			t.Errorf("Expected username 'getuser', got '%s'", user.Username)
		}
		if user.Email != "get@example.com" {
			t.Errorf("Expected email 'get@example.com', got '%s'", user.Email)
		}
		if user.CreatedAt.IsZero() {
			t.Error("Expected CreatedAt to be set")
		}
		if user.UpdatedAt.IsZero() {
			t.Error("Expected UpdatedAt to be set")
		}
	})

	t.Run("List users", func(t *testing.T) {
		// Create some additional users
		_, err := repo.CreateUser(ctx, "listuser1", "list1@example.com", testPassword)
		if err != nil {
			t.Fatalf("Failed to create test user 1: %v", err)
		}
		_, err = repo.CreateUser(ctx, "listuser2", "list2@example.com", testPassword)
		if err != nil {
			t.Fatalf("Failed to create test user 2: %v", err)
		}

		// Get the list of users
		users, err := repo.ListUsers(ctx)
		if err != nil {
			t.Fatalf("Failed to list users: %v", err)
		}

		// We should have at least 3 users
		if len(users) < 3 {
			t.Errorf("Expected at least 3 users, got %d", len(users))
		}

		// Check if our created users are in the list
		var found1, found2 bool
		for _, user := range users {
			if user.Username == "listuser1" && user.Email == "list1@example.com" {
				found1 = true
			}
			if user.Username == "listuser2" && user.Email == "list2@example.com" {
				found2 = true
			}
		}

		if !found1 {
			t.Error("listuser1 not found in user list")
		}
		if !found2 {
			t.Error("listuser2 not found in user list")
		}
	})

	t.Run("Update user", func(t *testing.T) {
		// First create a user to update
		id, err := repo.CreateUser(ctx, "updateuser", "update@example.com", testPassword)
		if err != nil {
			t.Fatalf("Failed to create user for update test: %v", err)
		}

		// Get the user to capture initial timestamps
		originalUser, err := repo.GetUser(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get original user: %v", err)
		}

		// Wait a moment to ensure timestamp changes
		time.Sleep(10 * time.Millisecond)

		// Update the user
		updatedUsername := "updateuser_new"
		updatedEmail := "update_new@example.com"
		err = repo.UpdateUser(ctx, id, updatedUsername, updatedEmail)
		if err != nil {
			t.Fatalf("Failed to update user: %v", err)
		}

		// Get the updated user
		updatedUser, err := repo.GetUser(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get updated user: %v", err)
		}

		// Validate updated data
		if updatedUser.Username != updatedUsername {
			t.Errorf("Expected username '%s', got '%s'", updatedUsername, updatedUser.Username)
		}
		if updatedUser.Email != updatedEmail {
			t.Errorf("Expected email '%s', got '%s'", updatedEmail, updatedUser.Email)
		}
		if !updatedUser.UpdatedAt.After(originalUser.UpdatedAt) {
			t.Error("Expected UpdatedAt to be updated")
		}
	})

	t.Run("Delete user", func(t *testing.T) {
		// First create a user to delete
		id, err := repo.CreateUser(ctx, "deleteuser", "delete@example.com", testPassword)
		if err != nil {
			t.Fatalf("Failed to create user for delete test: %v", err)
		}

		// Delete the user
		err = repo.DeleteUser(ctx, id)
		if err != nil {
			t.Fatalf("Failed to delete user: %v", err)
		}

		// Try to get the deleted user
		_, err = repo.GetUser(ctx, id)
		if err == nil {
			t.Error("Expected error when getting deleted user, got nil")
		}
	})

	t.Run("Get non-existent user", func(t *testing.T) {
		_, err := repo.GetUser(ctx, 9999)
		if err == nil {
			t.Error("Expected error when getting non-existent user, got nil")
		}
	})

	t.Run("Update non-existent user", func(t *testing.T) {
		err := repo.UpdateUser(ctx, 9999, "nonexistent", "nonexistent@example.com")
		if err == nil {
			t.Error("Expected error when updating non-existent user, got nil")
		}
	})

	t.Run("Delete non-existent user", func(t *testing.T) {
		err := repo.DeleteUser(ctx, 9999)
		if err == nil {
			t.Error("Expected error when deleting non-existent user, got nil")
		}
	})
}
