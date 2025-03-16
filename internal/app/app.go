package app

import (
	"database/sql"
	"html/template"
	"mordezzanV4/internal/controllers"
	"mordezzanV4/internal/logger"
	"mordezzanV4/internal/middleware"
	"mordezzanV4/internal/repositories"
	"mordezzanV4/internal/services"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	DB                    *sql.DB
	UserRepository        repositories.UserRepository
	CharacterRepository   repositories.CharacterRepository
	SpellRepository       repositories.SpellRepository
	ArmorRepository       repositories.ArmorRepository
	WeaponRepository      repositories.WeaponRepository
	EquipmentRepository   repositories.EquipmentRepository
	ShieldRepository      repositories.ShieldRepository
	PotionRepository      repositories.PotionRepository
	MagicItemRepository   repositories.MagicItemRepository
	RingRepository        repositories.RingRepository
	AmmoRepository        repositories.AmmoRepository
	SpellScrollRepository repositories.SpellScrollRepository
	ContainerRepository   repositories.ContainerRepository
	TreasureRepository    repositories.TreasureRepository
	InventoryRepository   repositories.InventoryRepository

	UserController        *controllers.UserController
	CharacterController   *controllers.CharacterController
	SpellController       *controllers.SpellController
	ArmorController       *controllers.ArmorController
	WeaponController      *controllers.WeaponController
	EquipmentController   *controllers.EquipmentController
	ShieldController      *controllers.ShieldController
	PotionController      *controllers.PotionController
	MagicItemController   *controllers.MagicItemController
	RingController        *controllers.RingController
	AmmoController        *controllers.AmmoController
	SpellScrollController *controllers.SpellScrollController
	ContainerController   *controllers.ContainerController
	AuthController        *controllers.AuthController
	TreasureController    *controllers.TreasureController
	InventoryController   *controllers.InventoryController

	Templates *template.Template
	JWTSecret string
}

func NewApp(dbPath string) (*App, error) {
	logger.Debug("Opening database connection to %s", dbPath)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		logger.Error("Failed to open database: %v", err)
		return nil, err
	}
	if err := db.Ping(); err != nil {
		logger.Error("Failed to ping database: %v", err)
		return nil, err
	}
	logger.Debug("Database connection established successfully")
	tmplPath := filepath.Join("web", "templates", "*.html")
	logger.Debug("Loading templates from %s", tmplPath)

	// Create a template with function map
	tmpl := template.New("").Funcs(template.FuncMap{
		"lower": strings.ToLower,
		"add":   func(a, b int) int { return a + b },
	})

	// Parse templates
	tmpl, err = tmpl.ParseGlob(tmplPath)
	if err != nil {
		logger.Error("Failed to parse templates: %v", err)
		return nil, err
	}
	logger.Debug("Templates loaded successfully")

	// Read JWT secret from environment or use a default for development
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		logger.Warning("JWT_SECRET not set, using default secret for development")
		jwtSecret = "mordezzan_development_secret_key_not_for_production"
	}

	userRepo := repositories.NewSQLCUserRepository(db)
	characterRepo := repositories.NewSQLCCharacterRepository(db)
	spellRepo := repositories.NewSQLCSpellRepository(db)
	armorRepo := repositories.NewSQLCArmorRepository(db)
	weaponRepo := repositories.NewSQLCWeaponRepository(db)
	equipmentRepo := repositories.NewSQLCEquipmentRepository(db)
	shieldRepo := repositories.NewSQLCShieldRepository(db)
	potionRepo := repositories.NewSQLCPotionRepository(db)
	magicItemRepo := repositories.NewSQLCMagicItemRepository(db)
	ringRepo := repositories.NewSQLCRingRepository(db)
	ammoRepo := repositories.NewSQLCAmmoRepository(db)
	spellScrollRepo := repositories.NewSQLCSpellScrollRepository(db)
	containerRepo := repositories.NewSQLCContainerRepository(db)
	treasureRepo := repositories.NewSQLCTreasureRepository(db)
	inventoryRepo := repositories.NewSQLCInventoryRepository(db)
	fighterDataRepo := repositories.NewSQLCFighterDataRepository(db)

	fighterService := services.NewFighterService(fighterDataRepo)

	userController := controllers.NewUserController(userRepo, tmpl)
	characterController := controllers.NewCharacterController(characterRepo, userRepo, fighterService, tmpl)
	spellController := controllers.NewSpellController(spellRepo, characterRepo, tmpl)
	armorController := controllers.NewArmorController(armorRepo, tmpl)
	weaponController := controllers.NewWeaponController(weaponRepo, tmpl)
	equipmentController := controllers.NewEquipmentController(equipmentRepo, tmpl)
	shieldController := controllers.NewShieldController(shieldRepo, tmpl)
	potionController := controllers.NewPotionController(potionRepo, tmpl)
	magicItemController := controllers.NewMagicItemController(magicItemRepo, tmpl)
	ringController := controllers.NewRingController(ringRepo, tmpl)
	ammoController := controllers.NewAmmoController(ammoRepo, tmpl)
	spellScrollController := controllers.NewSpellScrollController(spellScrollRepo, spellRepo, tmpl)
	containerController := controllers.NewContainerController(containerRepo, tmpl)
	authController := controllers.NewAuthController(userRepo, tmpl, jwtSecret)
	treasureController := controllers.NewTreasureController(treasureRepo, characterRepo, tmpl)
	inventoryController := controllers.NewInventoryController(
		inventoryRepo,
		characterRepo,
		weaponRepo,
		armorRepo,
		shieldRepo,
		potionRepo,
		magicItemRepo,
		ringRepo,
		ammoRepo,
		spellScrollRepo,
		containerRepo,
		equipmentRepo,
		treasureRepo,
		tmpl,
	)

	logger.Info("Application initialized successfully")

	return &App{
		DB:                    db,
		UserRepository:        userRepo,
		CharacterRepository:   characterRepo,
		SpellRepository:       spellRepo,
		ArmorRepository:       armorRepo,
		WeaponRepository:      weaponRepo,
		EquipmentRepository:   equipmentRepo,
		ShieldRepository:      shieldRepo,
		PotionRepository:      potionRepo,
		MagicItemRepository:   magicItemRepo,
		RingRepository:        ringRepo,
		AmmoRepository:        ammoRepo,
		SpellScrollRepository: spellScrollRepo,
		ContainerRepository:   containerRepo,
		TreasureRepository:    treasureRepo,
		InventoryRepository:   inventoryRepo,

		UserController:        userController,
		CharacterController:   characterController,
		SpellController:       spellController,
		ArmorController:       armorController,
		WeaponController:      weaponController,
		EquipmentController:   equipmentController,
		ShieldController:      shieldController,
		PotionController:      potionController,
		MagicItemController:   magicItemController,
		RingController:        ringController,
		AmmoController:        ammoController,
		SpellScrollController: spellScrollController,
		ContainerController:   containerController,
		AuthController:        authController,
		TreasureController:    treasureController,
		InventoryController:   inventoryController,

		Templates: tmpl,
		JWTSecret: jwtSecret,
	}, nil
}

func (a *App) SetupRoutes() http.Handler {
	logger.Debug("Setting up application routes")
	r := chi.NewRouter()

	// CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))).ServeHTTP)

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Home and dashboard routes (rendered HTML pages)
	r.Get("/", a.CharacterController.RenderDashboard)
	r.Get("/characters/create", a.CharacterController.RenderCreateForm)
	r.Get("/characters/view/{id}", a.CharacterController.RenderCharacterDetail)
	r.Get("/characters/{id}/edit", a.CharacterController.RenderEditForm)

	// Authentication routes
	r.Route("/auth", func(r chi.Router) {
		r.Get("/login-page", a.AuthController.RenderLoginPage)
		r.Get("/register-page", a.AuthController.RenderRegisterPage)
		r.Post("/login", a.AuthController.Login)
	})

	// API routes for data
	r.Route("/api", func(r chi.Router) {
		// User routes
		r.Route("/users", func(r chi.Router) {
			r.Get("/", a.UserController.ListUsers)
			r.Post("/", a.UserController.CreateUser)
			r.Get("/{id}", a.UserController.GetUser)
			r.Put("/{id}", a.UserController.UpdateUser)
			r.Delete("/{id}", a.UserController.DeleteUser)
			r.Get("/{id}/characters", a.CharacterController.GetCharactersByUser)
		})

		// Character routes
		r.Route("/characters", func(r chi.Router) {
			r.Get("/", a.CharacterController.ListCharacters)
			r.Post("/", a.CharacterController.CreateCharacter)
			r.Get("/{id}", a.CharacterController.GetCharacter)
			r.Put("/{id}", a.CharacterController.UpdateCharacter)
			r.Delete("/{id}", a.CharacterController.DeleteCharacter)
			r.Get("/{id}/spells", a.SpellController.GetSpellsByCharacter)
			r.Patch("/{id}/hp", a.CharacterController.UpdateCharacterHP)
			r.Patch("/{id}/xp", a.CharacterController.UpdateCharacterXP)
		})

		r.Route("/spells", func(r chi.Router) {
			r.Get("/", a.SpellController.ListSpells)
			r.Post("/", a.SpellController.CreateSpell)
			r.Get("/{id}", a.SpellController.GetSpell)
			r.Put("/{id}", a.SpellController.UpdateSpell)
			r.Delete("/{id}", a.SpellController.DeleteSpell)
		})

		r.Route("/armors", func(r chi.Router) {
			r.Get("/", a.ArmorController.ListArmors)
			r.Post("/", a.ArmorController.CreateArmor)
			r.Get("/{id}", a.ArmorController.GetArmor)
			r.Put("/{id}", a.ArmorController.UpdateArmor)
			r.Delete("/{id}", a.ArmorController.DeleteArmor)
		})

		r.Route("/weapons", func(r chi.Router) {
			r.Get("/", a.WeaponController.ListWeapons)
			r.Post("/", a.WeaponController.CreateWeapon)
			r.Get("/{id}", a.WeaponController.GetWeapon)
			r.Put("/{id}", a.WeaponController.UpdateWeapon)
			r.Delete("/{id}", a.WeaponController.DeleteWeapon)
		})

		r.Route("/equipment", func(r chi.Router) {
			r.Get("/", a.EquipmentController.ListEquipment)
			r.Post("/", a.EquipmentController.CreateEquipment)
			r.Get("/{id}", a.EquipmentController.GetEquipment)
			r.Put("/{id}", a.EquipmentController.UpdateEquipment)
			r.Delete("/{id}", a.EquipmentController.DeleteEquipment)
		})

		r.Route("/shields", func(r chi.Router) {
			r.Get("/", a.ShieldController.ListShields)
			r.Post("/", a.ShieldController.CreateShield)
			r.Get("/{id}", a.ShieldController.GetShield)
			r.Put("/{id}", a.ShieldController.UpdateShield)
			r.Delete("/{id}", a.ShieldController.DeleteShield)
		})

		r.Route("/potions", func(r chi.Router) {
			r.Get("/", a.PotionController.ListPotions)
			r.Post("/", a.PotionController.CreatePotion)
			r.Get("/{id}", a.PotionController.GetPotion)
			r.Put("/{id}", a.PotionController.UpdatePotion)
			r.Delete("/{id}", a.PotionController.DeletePotion)
		})

		r.Route("/magic-items", func(r chi.Router) {
			r.Get("/", a.MagicItemController.ListMagicItems)
			r.Post("/", a.MagicItemController.CreateMagicItem)
			r.Get("/{id}", a.MagicItemController.GetMagicItem)
			r.Put("/{id}", a.MagicItemController.UpdateMagicItem)
			r.Delete("/{id}", a.MagicItemController.DeleteMagicItem)
		})

		r.Route("/rings", func(r chi.Router) {
			r.Get("/", a.RingController.ListRings)
			r.Post("/", a.RingController.CreateRing)
			r.Get("/{id}", a.RingController.GetRing)
			r.Put("/{id}", a.RingController.UpdateRing)
			r.Delete("/{id}", a.RingController.DeleteRing)
		})

		r.Route("/ammo", func(r chi.Router) {
			r.Get("/", a.AmmoController.ListAmmo)
			r.Post("/", a.AmmoController.CreateAmmo)
			r.Get("/{id}", a.AmmoController.GetAmmo)
			r.Put("/{id}", a.AmmoController.UpdateAmmo)
			r.Delete("/{id}", a.AmmoController.DeleteAmmo)
		})

		r.Route("/spell-scrolls", func(r chi.Router) {
			r.Get("/", a.SpellScrollController.ListSpellScrolls)
			r.Post("/", a.SpellScrollController.CreateSpellScroll)
			r.Get("/{id}", a.SpellScrollController.GetSpellScroll)
			r.Put("/{id}", a.SpellScrollController.UpdateSpellScroll)
			r.Delete("/{id}", a.SpellScrollController.DeleteSpellScroll)
		})

		r.Route("/containers", func(r chi.Router) {
			r.Get("/", a.ContainerController.ListContainers)
			r.Post("/", a.ContainerController.CreateContainer)
			r.Get("/{id}", a.ContainerController.GetContainer)
			r.Put("/{id}", a.ContainerController.UpdateContainer)
			r.Delete("/{id}", a.ContainerController.DeleteContainer)
		})

		r.Route("/treasures", func(r chi.Router) {
			r.Get("/", a.TreasureController.ListTreasures)
			r.Post("/", a.TreasureController.CreateTreasure)
			r.Get("/{id}", a.TreasureController.GetTreasure)
			r.Put("/{id}", a.TreasureController.UpdateTreasure)
			r.Delete("/{id}", a.TreasureController.DeleteTreasure)
			r.Get("/character/{characterId}", a.TreasureController.GetTreasureByCharacter)
		})

		r.Route("/inventories", func(r chi.Router) {
			r.Get("/", a.InventoryController.ListInventories)
			r.Post("/", a.InventoryController.CreateInventory)
			r.Get("/{id}", a.InventoryController.GetInventory)
			r.Put("/{id}", a.InventoryController.UpdateInventory)
			r.Delete("/{id}", a.InventoryController.DeleteInventory)
			r.Post("/{id}/items", a.InventoryController.AddInventoryItem)
			r.Get("/{id}/items/{itemId}", a.InventoryController.GetInventoryItem)
			r.Put("/{id}/items/{itemId}", a.InventoryController.UpdateInventoryItem)
			r.Delete("/{id}/items/{itemId}", a.InventoryController.RemoveInventoryItem)
			r.Get("/character/{characterId}", a.InventoryController.GetInventoryByCharacter)
		})
	})

	// Configure JWT authentication
	authConfig := middleware.AuthConfig{
		JWTSecret: a.JWTSecret,
		Issuer:    "mordezzanV4",
	}

	// Apply JWT auth middleware to protected routes
	handler := middleware.RecoveryMiddleware(
		middleware.LoggingMiddleware(
			middleware.JWTAuthMiddleware(authConfig)(r),
		),
	)

	logger.Info("Routes set up successfully")
	return handler
}

func (a *App) Shutdown() {
	if a.DB != nil {
		logger.Info("Closing database connection...")
		a.DB.Close()
	}
}
