package integration

import (
	"database/sql"
	"log"
	"mordezzanV4/internal/app"
	"mordezzanV4/internal/logger"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// TestMain sets up the test environment for all integration tests
func TestMain(m *testing.M) {
	log.Println("Initializing logger for all integration tests...")
	logger.Init(logger.Config{
		LogLevel:         logger.LogLevelDebug,
		IncludeTimestamp: true,
		IncludeFileLine:  true,
		Output:           os.Stdout,
	})
	exitCode := m.Run()
	os.Exit(exitCode)
}

// init initializes the logger for this package
func init() {
	logger.Init(logger.Config{
		LogLevel:         logger.LogLevelInfo,
		IncludeTimestamp: true,
		IncludeFileLine:  true,
		Output:           os.Stdout,
	})
}

func setupTestEnvironment(t *testing.T) (*app.App, *httptest.Server, *TestUser, func()) {
	testApp, appCleanup := setupTestAppWithFullSchema(t)

	handler := testApp.SetupRoutes()
	server := httptest.NewServer(handler)

	// Create authenticated user
	testUser := CreateTestUserWithAuth(t, server)

	cleanup := func() {
		server.Close()
		appCleanup()
	}

	return testApp, server, testUser, cleanup
}

// setupTestApp creates a test app instance with a temporary database
func setupTestApp(t *testing.T) (*app.App, func()) {
	tempDir, err := os.MkdirTemp("", "test_templates")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create user template
	templatePath := filepath.Join(tempDir, "user.html")
	templateContent := `<!DOCTYPE html><html><body>{{.Username}}</body></html>`
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("Failed to write template: %v", err)
	}

	// Create templates directory if it doesn't exist
	if err := os.MkdirAll("web/templates", 0755); err != nil {
		t.Fatalf("Failed to create template dir: %v", err)
	}

	// Create symlink or copy the template file
	if err := os.Symlink(templatePath, "web/templates/user.html"); err != nil {
		input, _ := os.ReadFile(templatePath)
		os.WriteFile("web/templates/user.html", input, 0644)
	}

	// Create a test database
	dbFile := "./integration_test.db"
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Create users table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
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

	db.Close()

	// Initialize the test app
	testApp, err := app.NewApp(dbFile)
	if err != nil {
		t.Fatalf("Failed to initialize test app: %v", err)
	}

	// Return cleanup function
	return testApp, func() {
		testApp.Shutdown()
		os.Remove(dbFile)
		os.RemoveAll(tempDir)
		os.Remove("web/templates/user.html")
	}
}

func setupTestAppWithFullSchema(t *testing.T) (*app.App, func()) {
	tempDir, err := os.MkdirTemp("", "test_templates")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	templatePath := filepath.Join(tempDir, "user.html")
	templateContent := `<!DOCTYPE html><html><body>{{.Username}}</body></html>`
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("Failed to write template: %v", err)
	}
	if err := os.MkdirAll("web/templates", 0755); err != nil {
		t.Fatalf("Failed to create template dir: %v", err)
	}
	if err := os.Symlink(templatePath, "web/templates/user.html"); err != nil {
		input, _ := os.ReadFile(templatePath)
		os.WriteFile("web/templates/user.html", input, 0644)
	}

	characterTemplatePath := filepath.Join(tempDir, "character.html")
	characterTemplateContent := `<!DOCTYPE html><html><body>{{.Name}}</body></html>`
	if err := os.WriteFile(characterTemplatePath, []byte(characterTemplateContent), 0644); err != nil {
		t.Fatalf("Failed to write character template: %v", err)
	}
	if err := os.Symlink(characterTemplatePath, "web/templates/character.html"); err != nil {
		input, _ := os.ReadFile(characterTemplatePath)
		os.WriteFile("web/templates/character.html", input, 0644)
	}

	spellTemplatePath := filepath.Join(tempDir, "spell.html")
	spellTemplateContent := `<!DOCTYPE html><html><body>{{.Name}}</body></html>`
	if err := os.WriteFile(spellTemplatePath, []byte(spellTemplateContent), 0644); err != nil {
		t.Fatalf("Failed to write spell template: %v", err)
	}
	if err := os.Symlink(spellTemplatePath, "web/templates/spell.html"); err != nil {
		input, _ := os.ReadFile(spellTemplatePath)
		os.WriteFile("web/templates/spell.html", input, 0644)
	}

	armorTemplatePath := filepath.Join(tempDir, "armor.html")
	armorTemplateContent := `<!DOCTYPE html><html><body>{{.Name}}</body></html>`
	if err := os.WriteFile(armorTemplatePath, []byte(armorTemplateContent), 0644); err != nil {
		t.Fatalf("Failed to write armor template: %v", err)
	}
	if err := os.Symlink(armorTemplatePath, "web/templates/armor.html"); err != nil {
		input, _ := os.ReadFile(armorTemplatePath)
		os.WriteFile("web/templates/armor.html", input, 0644)
	}

	weaponTemplatePath := filepath.Join(tempDir, "weapon.html")
	weaponTemplateContent := `<!DOCTYPE html><html><body>{{.Name}}</body></html>`
	if err := os.WriteFile(weaponTemplatePath, []byte(weaponTemplateContent), 0644); err != nil {
		t.Fatalf("Failed to write weapon template: %v", err)
	}
	if err := os.Symlink(weaponTemplatePath, "web/templates/weapon.html"); err != nil {
		input, _ := os.ReadFile(weaponTemplatePath)
		os.WriteFile("web/templates/weapon.html", input, 0644)
	}

	shieldTemplatePath := filepath.Join(tempDir, "shield.html")
	shieldTemplateContent := `<!DOCTYPE html><html><body>{{.Name}}</body></html>`
	if err := os.WriteFile(shieldTemplatePath, []byte(shieldTemplateContent), 0644); err != nil {
		t.Fatalf("Failed to write shield template: %v", err)
	}
	if err := os.Symlink(shieldTemplatePath, "web/templates/shield.html"); err != nil {
		input, _ := os.ReadFile(shieldTemplatePath)
		os.WriteFile("web/templates/shield.html", input, 0644)
	}

	equipmentTemplatePath := filepath.Join(tempDir, "equipment.html")
	equipmentTemplateContent := `<!DOCTYPE html><html><body>{{.Name}}</body></html>`
	if err := os.WriteFile(equipmentTemplatePath, []byte(equipmentTemplateContent), 0644); err != nil {
		t.Fatalf("Failed to write equipment template: %v", err)
	}
	if err := os.Symlink(equipmentTemplatePath, "web/templates/equipment.html"); err != nil {
		input, _ := os.ReadFile(equipmentTemplatePath)
		os.WriteFile("web/templates/equipment.html", input, 0644)
	}

	magicItemTemplatePath := filepath.Join(tempDir, "magic_item.html")
	magicItemTemplateContent := `<!DOCTYPE html><html><body>{{.Name}}</body></html>`
	if err := os.WriteFile(magicItemTemplatePath, []byte(magicItemTemplateContent), 0644); err != nil {
		t.Fatalf("Failed to write magic item template: %v", err)
	}
	if err := os.Symlink(magicItemTemplatePath, "web/templates/magic_item.html"); err != nil {
		input, _ := os.ReadFile(magicItemTemplatePath)
		os.WriteFile("web/templates/magic_item.html", input, 0644)
	}

	ringTemplatePath := filepath.Join(tempDir, "ring.html")
	ringTemplateContent := `<!DOCTYPE html><html><body>{{.Name}}</body></html>`
	if err := os.WriteFile(ringTemplatePath, []byte(ringTemplateContent), 0644); err != nil {
		t.Fatalf("Failed to write ring template: %v", err)
	}
	if err := os.Symlink(ringTemplatePath, "web/templates/ring.html"); err != nil {
		input, _ := os.ReadFile(ringTemplatePath)
		os.WriteFile("web/templates/ring.html", input, 0644)
	}

	ammoTemplatePath := filepath.Join(tempDir, "ammo.html")
	ammoTemplateContent := `<!DOCTYPE html><html><body>{{.Name}}</body></html>`
	if err := os.WriteFile(ammoTemplatePath, []byte(ammoTemplateContent), 0644); err != nil {
		t.Fatalf("Failed to write ammo template: %v", err)
	}
	if err := os.Symlink(ammoTemplatePath, "web/templates/ammo.html"); err != nil {
		input, _ := os.ReadFile(ammoTemplatePath)
		os.WriteFile("web/templates/ammo.html", input, 0644)
	}

	spellScrollTemplatePath := filepath.Join(tempDir, "spell_scroll.html")
	spellScrollTemplateContent := `<!DOCTYPE html><html><body>{{.SpellName}}</body></html>`
	if err := os.WriteFile(spellScrollTemplatePath, []byte(spellScrollTemplateContent), 0644); err != nil {
		t.Fatalf("Failed to write spell scroll template: %v", err)
	}
	if err := os.Symlink(spellScrollTemplatePath, "web/templates/spell_scroll.html"); err != nil {
		input, _ := os.ReadFile(spellScrollTemplatePath)
		os.WriteFile("web/templates/spell_scroll.html", input, 0644)
	}

	containerTemplatePath := filepath.Join(tempDir, "container.html")
	containerTemplateContent := `<!DOCTYPE html><html><body>{{.Name}}</body></html>`
	if err := os.WriteFile(containerTemplatePath, []byte(containerTemplateContent), 0644); err != nil {
		t.Fatalf("Failed to write container template: %v", err)
	}
	if err := os.Symlink(containerTemplatePath, "web/templates/container.html"); err != nil {
		input, _ := os.ReadFile(containerTemplatePath)
		os.WriteFile("web/templates/container.html", input, 0644)
	}

	treasureTemplatePath := filepath.Join(tempDir, "treasure.html")
	treasureTemplateContent := `<!DOCTYPE html><html><body>Character ID: {{.CharacterID}}, Gold: {{.GoldCoins}}</body></html>`
	if err := os.WriteFile(treasureTemplatePath, []byte(treasureTemplateContent), 0644); err != nil {
		t.Fatalf("Failed to write treasure template: %v", err)
	}
	if err := os.Symlink(treasureTemplatePath, "web/templates/treasure.html"); err != nil {
		input, _ := os.ReadFile(treasureTemplatePath)
		os.WriteFile("web/templates/treasure.html", input, 0644)
	}

	dbFile := "./character_integration_test.db"
	log.Printf("Setting up test database at %s", dbFile)
	os.Remove(dbFile)
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	log.Println("Creating users table...")
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT NOT NULL UNIQUE,
            email TEXT NOT NULL UNIQUE,
            password_hash TEXT NOT NULL,
            created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
        );
    `)
	if err != nil {
		t.Fatalf("Failed to create users table: %v", err)
	}

	log.Println("Creating characters table...")
	_, err = db.Exec(`
		CREATE TABLE characters (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			class TEXT NOT NULL DEFAULT 'Fighter',
			level INTEGER NOT NULL DEFAULT 1,
			strength INTEGER NOT NULL DEFAULT 10,
			dexterity INTEGER NOT NULL DEFAULT 10,
			constitution INTEGER NOT NULL DEFAULT 10,
			wisdom INTEGER NOT NULL DEFAULT 10,
			intelligence INTEGER NOT NULL DEFAULT 10,
			charisma INTEGER NOT NULL DEFAULT 10,
			hit_points INTEGER NOT NULL DEFAULT 10,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		db.Close()
		t.Fatalf("Failed to create characters table: %v", err)
	}

	log.Println("Creating spells table...")
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS spells (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            character_id INTEGER NOT NULL,
            name TEXT NOT NULL,
            mag_level INTEGER NOT NULL DEFAULT 0,
            cry_level INTEGER NOT NULL DEFAULT 0,
            ill_level INTEGER NOT NULL DEFAULT 0,
            nec_level INTEGER NOT NULL DEFAULT 0,
            pyr_level INTEGER NOT NULL DEFAULT 0,
            wch_level INTEGER NOT NULL DEFAULT 0,
            clr_level INTEGER NOT NULL DEFAULT 0,
            drd_level INTEGER NOT NULL DEFAULT 0,
            range TEXT NOT NULL,
            duration TEXT NOT NULL,
            area_of_effect TEXT,
            components TEXT,
            description TEXT NOT NULL,
            created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (character_id) REFERENCES characters (id) ON DELETE CASCADE
        );
    `)
	if err != nil {
		db.Close()
		t.Fatalf("Failed to create spells table: %v", err)
	}

	log.Println("Creating armors table...")
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS armors (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL UNIQUE,
            armor_type TEXT NOT NULL,
            ac INTEGER NOT NULL,
            cost REAL NOT NULL,
            damage_reduction INTEGER NOT NULL DEFAULT 0,
            weight INTEGER NOT NULL,
            weight_class TEXT NOT NULL,
            movement_rate INTEGER NOT NULL,
            created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
        );
    `)
	if err != nil {
		db.Close()
		t.Fatalf("Failed to create armors table: %v", err)
	}

	log.Println("Creating weapons table...")
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS weapons (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            category TEXT NOT NULL,
            weapon_class INTEGER NOT NULL DEFAULT 1,
            cost REAL NOT NULL,
            weight INTEGER NOT NULL,
            range_short INTEGER,
            range_medium INTEGER,
            range_long INTEGER,
            rate_of_fire TEXT,
            damage TEXT NOT NULL,
            damage_two_handed TEXT,
            properties TEXT,
            created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
        );
    `)
	if err != nil {
		db.Close()
		t.Fatalf("Failed to create weapons table: %v", err)
	}

	log.Println("Creating equipment table...")
	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS equipment (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        description TEXT NOT NULL,
        cost REAL NOT NULL,
        weight INTEGER NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    	);
	`)
	if err != nil {
		db.Close()
		t.Fatalf("Failed to create equipment table: %v", err)
	}

	log.Println("Creating shields table...")
	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS shields (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        cost REAL NOT NULL,
        weight INTEGER NOT NULL,
        defense_modifier INTEGER NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
	`)
	if err != nil {
		db.Close()
		t.Fatalf("Failed to create shields table: %v", err)
	}
	log.Println("Creating potions table...")
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS potions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		uses INTEGER NOT NULL DEFAULT 1,
		weight INTEGER NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		db.Close()
		t.Fatalf("Failed to create potions table: %v", err)
	}

	log.Println("Creating magic_items table...")
	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS magic_items (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        item_type TEXT NOT NULL,
        description TEXT NOT NULL,
        charges INTEGER,
        cost REAL NOT NULL,
        weight INTEGER NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
	`)
	if err != nil {
		db.Close()
		t.Fatalf("Failed to create magic_items table: %v", err)
	}

	log.Println("Creating rings table...")
	_, err = db.Exec(`
	CREATE TABLE rings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		cost REAL NOT NULL,
		weight INTEGER NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	`)
	if err != nil {
		db.Close()
		t.Fatalf("Failed to create rings table: %v", err)
	}

	log.Println("Creating ammo table...")
	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS ammo (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        cost REAL NOT NULL,
        weight INTEGER NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
	`)
	if err != nil {
		db.Close()
		t.Fatalf("Failed to create ammo table: %v", err)
	}

	log.Println("Creating spell_scrolls table...")
	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS spell_scrolls (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        spell_id INTEGER NOT NULL,
        casting_level INTEGER NOT NULL DEFAULT 1,
        cost REAL NOT NULL,
        weight INTEGER NOT NULL,
        description TEXT,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (spell_id) REFERENCES spells (id) ON DELETE CASCADE
    );
	`)
	if err != nil {
		db.Close()
		t.Fatalf("Failed to create spell_scrolls table: %v", err)
	}

	log.Println("Creating containers table...")
	_, err = db.Exec(`
	CREATE TABLE containers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		max_weight INTEGER NOT NULL,
		allowed_items TEXT NOT NULL,
		cost REAL NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	`)
	if err != nil {
		db.Close()
		t.Fatalf("Failed to create containers table: %v", err)
	}

	log.Println("Creating treasures table...")
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS treasures (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		character_id INTEGER NULL,
		platinum_coins INTEGER NOT NULL DEFAULT 0,
		gold_coins INTEGER NOT NULL DEFAULT 0,
		electrum_coins INTEGER NOT NULL DEFAULT 0,
		silver_coins INTEGER NOT NULL DEFAULT 0,
		copper_coins INTEGER NOT NULL DEFAULT 0,
		gems TEXT NULL,
		art_objects TEXT NULL,
		other_valuables TEXT NULL,
		total_value_gold REAL NOT NULL DEFAULT 0,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (character_id) REFERENCES characters (id) ON DELETE SET NULL
		);

	`)
	if err != nil {
		db.Close()
		t.Fatalf("Failed to create treasures table: %v", err)
	}

	log.Println("Creating inventories table...")
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS inventories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		character_id INTEGER NOT NULL,
		max_weight REAL NOT NULL DEFAULT 100.0,
		current_weight REAL NOT NULL DEFAULT 0.0,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (character_id) REFERENCES characters (id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		db.Close()
		t.Fatalf("Failed to create inventories table: %v", err)
	}

	log.Println("Creating inventory_items table...")
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS inventory_items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		inventory_id INTEGER NOT NULL,
		item_type TEXT NOT NULL,
		item_id INTEGER NOT NULL,
		quantity INTEGER NOT NULL DEFAULT 1,
		is_equipped BOOLEAN NOT NULL DEFAULT 0,
		notes TEXT,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (inventory_id) REFERENCES inventories (id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		db.Close()
		t.Fatalf("Failed to create inventory_items table: %v", err)
	}

	db.Close()
	log.Println("Database tables created successfully")

	log.Println("Initializing test app...")
	testApp, err := app.NewApp(dbFile)
	if err != nil {
		t.Fatalf("Failed to initialize test app: %v", err)
	}

	return testApp, func() {
		log.Println("Cleaning up test resources...")
		testApp.Shutdown()
		os.Remove(dbFile)
		os.RemoveAll(tempDir)
		os.Remove("web/templates/user.html")
		os.Remove("web/templates/character.html")
		os.Remove("web/templates/spell.html")
		os.Remove("web/templates/armor.html")
		os.Remove("web/templates/weapon.html")
		os.Remove("web/templates/shield.html")
		os.Remove("web/templates/equipment.html")
		os.Remove("web/templates/magic_item.html")
		os.Remove("web/templates/ring.html")
		os.Remove("web/templates/ammo.html")
		os.Remove("web/templates/spell_scroll.html")
		os.Remove("web/templates/container.html")
		os.Remove("web/templates/treasure.html")
		log.Println("Cleanup completed")
	}
}
