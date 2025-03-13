# Go MVC Project

This is a Go web application following the MVC (Model-View-Controller) pattern.

## Technologies

- Go
- SQLite (Database)
- SQLC (Type-safe SQL queries)
- Goose (Database migrations)

## Getting Started

1. Install required tools:
   ```
   go install github.com/pressly/goose/v3/cmd/goose@latest
   go install github.com/kyleconroy/sqlc/cmd/sqlc@latest
   ```

2. Set up the project:
   ```
   make dev-setup
   ```

3. Run the application:
   ```
   make run
   ```

## Project Structure

- `cmd/server`: Application entry point
- `internal/models`: Domain models and business logic
- `internal/controllers`: Request handlers
- `internal/views`: Templates and presentation logic
- `internal/repositories`: Data access layer
- `internal/middleware`: HTTP middleware components
- `pkg`: Reusable packages
- `web`: Web assets (CSS, JS, images, templates)
