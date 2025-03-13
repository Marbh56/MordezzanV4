# Testing Guide

This document outlines the testing approach for the Go MVC Project.

## Test Structure

The project includes several levels of testing:

1. **Unit Tests**: Testing individual components in isolation
   - Repository tests: Test database operations using a test database
   - Controller tests: Test HTTP handlers using mocked repositories
   - Middleware tests: Test middleware behavior

2. **Integration Tests**: Testing the entire application flow
   - End-to-end tests that verify CRUD operations through the API

## Running Tests

You can run tests using make commands:

```bash
# Run all tests
make test-all

# Run only unit tests
make test-unit

# Run only integration tests
make test-integration

# Run tests with Go's standard test command
go test ./...
```

## Test Files

- `internal/repositories/user_repository_test.go`: Tests for the user repository
- `internal/controllers/user_controller_test.go`: Tests for the user controller
- `internal/middleware/auth_test.go`: Tests for the authentication middleware
- `integration_test/integration_test.go`: Integration tests for the entire application

## Writing New Tests

When adding new features, follow these guidelines for testing:

1. **Unit Tests**:
   - Test each component in isolation
   - Use mocks for dependencies
   - Focus on all edge cases

2. **Integration Tests**:
   - Test the complete flow
   - Use a temporary database
   - Clean up resources after tests

## Mock Testing

For controller tests, we use a mock user repository:

```go
// Example of a mock repository
type MockUserRepository struct {
    users map[int64]*models.User
    // ...
}

func (m *MockUserRepository) GetUser(ctx context.Context, id int64) (*models.User, error) {
    // Mock implementation
}
```

## Test Database

For tests that require a database:

1. Create a temporary SQLite database file
2. Set up the necessary schema
3. Run tests
4. Clean up the database file afterward

Example:
```go
func setupTestDB(t *testing.T) (*sql.DB, func()) {
    dbFile := "./test_db.db"
    db, err := sql.Open("sqlite3", dbFile)
    // ...setup schema...
    
    return db, func() {
        db.Close()
        os.Remove(dbFile)
    }
}
```

## Continuous Integration

In a CI environment, tests are run using:

```bash
make test-all
```

This ensures all new code passes both unit and integration tests before merging.