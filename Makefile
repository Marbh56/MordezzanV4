.PHONY: build run test migrate sqlc clean test-integration test-integration-simple

build:
	go build -o bin/server cmd/server/main.go

run: build
	./bin/server

test:
	go test ./...

test-unit:
	go test ./internal/repositories/... ./internal/controllers/... ./internal/middleware/... -v

test-integration:
	bash integration_test/test-run-script.sh

test-integration-simple:
	bash integration_test/test-run-script.sh --simple

test-all: test-unit test-integration

integration-test: test-integration

create-migration:
	@if [ -z "$(name)" ]; then \
		echo "Error: Migration name is required. Usage: make create-migration name=migration_name"; \
		exit 1; \
	fi
	@echo "Creating migration: $(name)"
	@timestamp=$$(date +%Y%m%d%H%M%S); \
	filename="internal/repositories/db/migrations/$${timestamp}_$(name).sql"; \
	echo "-- +goose Up\n-- SQL in this section is executed when the migration is applied\n\n\n-- +goose Down\n-- SQL in this section is executed when the migration is rolled back\n" > "$$filename"; \
	echo "Created migration file: $$filename"

migrate:
	goose -dir internal/repositories/db/migrations sqlite3 ./myproject.db up

sqlc:
	sqlc generate

clean:
	go clean

dev-setup: sqlc migrate build

repo:
	repomix --remove-comments --remove-empty-lines

repo-ignore:
	repomix --remove-comments --remove-empty-lines --ignore "**/*.log,**/*.db,**/bin/**,**/tmp/**,**/.git/**,**/node_modules/**,**/.DS_Store,**/*.sqlite,**/*.sqlite3"

repo-clean:
	repomix --remove-comments --remove-empty-lines \
		--ignore "**/*_test.go,**/integration_test/**,**/*test*.go,**/*.log,**/*.db,**/bin/**,**/tmp/**,**/.git/**,**/node_modules/**,**/.DS_Store,**/*.sqlite,**/*.sqlite3,**/test_logs/**,**/test_*.*,test-run-script.sh"

repo-frontend:
	repomix --remove-comments --remove-empty-lines \
		--ignore "**/*repository*.go,**/*.sql,**/*.sql.go,**/queries/**,**/*_test.go,**/db/**,**/migrations/**,**/sqlc/**,**/integration_test/**,**/*test*.go,**/*.log,**/*.db,**/bin/**,**/tmp/**,**/.git/**,**/node_modules/**,**/.DS_Store,**/*.sqlite,**/*.sqlite3,**/test_logs/**,**/test_*.*,test-run-script.sh,**/middleware/**,cmd/**,**/logger/**,**/errors/**,**/app/**,**/contextkeys/**"

repo-backend:
	repomix --remove-comments --remove-empty-lines \
		--ignore "**/*_test.go,**/integration_test/**,**/*test*.go,**/*.log,**/*.db,**/bin/**,**/tmp/**,**/.git/**,**/node_modules/**,**/.DS_Store,**/*.sqlite,**/*.sqlite3,**/test_logs/**,**/test_*.*,test-run-script.sh,**/web/static/**,**/web/templates/**,**/docs/**,README.md,**/migrations/*.sql,**/queries/*.sql,**/sqlc/models.go,**/sqlc/db.go,**/sqlc/querier.go,**/sqlc/*.sql.go,Makefile,sqlc.yaml,.gitignore,go.mod,**/ammo*.go,**/armor*.go,**/shield*.go,**/ring*.go,**/potion*.go,**/weapon*.go,**/equipment*.go,**/magic_item*.go,**/spell_scroll*.go,**/container*.go,**/spell*.go"

repo-dashboard:
	repomix --remove-comments --remove-empty-lines \
		--ignore "**/*repository*.go,**/*.sql,**/*.sql.go,**/queries/**,**/*_test.go,**/db/**,**/migrations/**,**/sqlc/**,**/integration_test/**,**/*test*.go,**/*.log,**/*.db,**/bin/**,**/tmp/**,**/.git/**,**/node_modules/**,**/.DS_Store,**/*.sqlite,**/*.sqlite3,**/test_logs/**,**/test_*.*,test-run-script.sh,**/ammo*.go,**/armor*.go,**/shield*.go,**/ring*.go,**/potion*.go,**/weapon*.go,**/equipment*.go,**/magic_item*.go,**/spell*.go,**/spell_scroll*.go,**/container*.go,**/treasure*.go,**/inventory*.go"

repo-fighter:
	repomix --remove-comments --remove-empty-lines \
		--ignore "**/ammo*,**/armor*,**/treasure*,**/weapon*,**/equipment*,**/inventory*,**/shield*,**/spell*,**/magic*,**/potion*,**/ring*,**/container*,*test*,**/*.log,**/*.db,**/bin/**,**/tmp/**,**/.git/**,**/node_modules/**,**/.DS_Store,**/*.sqlite,**/*.sqlite3,**/test_logs/**,**/auth/**,**/contextkeys/**,**/errors/**,**/logger/**,**/middleware/**"