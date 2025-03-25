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

tree:
	@if command -v tree > /dev/null; then \
		tree -I "node_modules|vendor|.git" --dirsfirst -F > project_structure.txt; \
		echo "Project structure saved to project_structure.txt"; \
	else \
		find . -type d -not -path "*/\.*" -not -path "*/node_modules/*" -not -path "*/vendor/*" | sort > project_structure.txt; \
		find . -type f -not -path "*/\.*" -not -path "*/node_modules/*" -not -path "*/vendor/*" | sort >> project_structure.txt; \
		echo "Project structure saved to project_structure.txt (using find command)"; \
	fi

repo-dashboard:
	repomix --remove-comments --remove-empty-lines \
		--ignore "**/controllers/ammo_controller.go,**/controllers/armor_controller.go,**/controllers/spell_controller.go,**/controllers/shield_controller.go,**/controllers/weapon_controller.go,**/controllers/equipment_controller.go,**/controllers/magic_items_controller.go,**/controllers/potion_controller.go,**/controllers/ring_controller.go,**/controllers/spell_scroll_controller.go,**/controllers/container_controller.go,**/controllers/treasure_controller.go,**/controllers/inventory_controller.go,**/controllers/spell_casting_controller.go,**/controllers/character_controller.go,**/models/ammo.go,**/models/armor.go,**/models/shield.go,**/models/ring.go,**/models/potion.go,**/models/weapon.go,**/models/equipment.go,**/models/magic_items.go,**/models/spell_scrolls.go,**/models/containers.go,**/models/treasure.go,**/models/spell*.go,**/models/class_data.go,**/models/encumbrance.go,**/models/inventory.go,**/models/character.go,**/repositories/ammo_repository.go,**/repositories/armor_repository.go,**/repositories/shield_repository.go,**/repositories/ring_repository.go,**/repositories/potion_repository.go,**/repositories/weapon_repository.go,**/repositories/equipment_repository.go,**/repositories/magic_items_repository.go,**/repositories/spell_scroll_repository.go,**/repositories/container_repository.go,**/repositories/treasure_repository.go,**/repositories/spell_repository.go,**/repositories/character_repository.go,**/repositories/class_repository.go,**/repositories/inventory_repository.go,**/repositories/spell_casting_repository.go,**/services/class_service.go,**/services/encumbrance_service.go,**/services/spell_service.go,**/queries/**,**/*_test.go,**/db/migrations/**,**/db/inserts/**,**/sqlc/**,**/integration_test/**,**/*.log,**/*.db,**/bin/**,**/tmp/**,**/.git/**,**/node_modules/**,**/.DS_Store,**/*.sqlite,**/*.sqlite3,**/test_logs/**,**/web/static/css/spells_tab.css,**/web/static/css/inventory_tab.css,**/web/static/js/spells_tab.js,**/web/static/js/inventory_tab.js,**/web/templates/ammo.html,**/web/templates/armor.html,**/web/templates/potion.html,**/web/templates/ring.html,**/web/templates/shield.html,**/web/templates/weapon.html,**/web/templates/magic_items.html,**/web/templates/containers.html,**/web/templates/equipment.html,**/web/templates/spell*.html,**/web/templates/inventory.html,**/web/templates/character*.html"

repo-character:
	repomix --remove-comments --remove-empty-lines \
		--ignore "**/controllers/ammo_controller.go,**/controllers/armor_controller.go,**/controllers/spell_controller.go,**/controllers/shield_controller.go,**/controllers/weapon_controller.go,**/controllers/equipment_controller.go,**/controllers/magic_items_controller.go,**/controllers/potion_controller.go,**/controllers/ring_controller.go,**/controllers/spell_scroll_controller.go,**/controllers/container_controller.go,**/controllers/treasure_controller.go,**/controllers/inventory_controller.go,**/controllers/spell_casting_controller.go,**/controllers/auth_controller.go,**/controllers/user_controller.go,**/controllers/base_controller.go,**/models/ammo.go,**/models/armor.go,**/models/shield.go,**/models/ring.go,**/models/potion.go,**/models/weapon.go,**/models/equipment.go,**/models/magic_items.go,**/models/spell_scrolls.go,**/models/containers.go,**/models/treasure.go,**/models/spell*.go,**/models/encumbrance.go,**/models/inventory.go,**/models/user.go,**/repositories/**,**/services/**,**/middleware/**,**/db/migrations/**,**/db/inserts/**,**/sqlc/**,**/integration_test/**,**/*.log,**/*.db,**/bin/**,**/tmp/**,**/.git/**,**/node_modules/**,**/.DS_Store,**/*.sqlite,**/*.sqlite3,**/test_logs/**,**/web/static/css/*tab.css,**/web/static/js/*tab.js,**/web/templates/ammo.html,**/web/templates/armor.html,**/web/templates/potion.html,**/web/templates/ring.html,**/web/templates/shield.html,**/web/templates/weapon.html,**/web/templates/magic_items.html,**/web/templates/containers.html,**/web/templates/equipment.html,**/web/templates/spell*.html,**/web/templates/inventory.html,**/web/templates/login.html,**/web/templates/register.html,**/web/templates/home.html,**/web/templates/dashboard.html"