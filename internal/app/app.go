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
	"path/filepath"
	"strings"
	"time"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	DB                      *sql.DB
	UserRepository          repositories.UserRepository
	CharacterRepository     repositories.CharacterRepository
	SpellRepository         repositories.SpellRepository
	ArmorRepository         repositories.ArmorRepository
	WeaponRepository        repositories.WeaponRepository
	EquipmentRepository     repositories.EquipmentRepository
	ShieldRepository        repositories.ShieldRepository
	PotionRepository        repositories.PotionRepository
	MagicItemRepository     repositories.MagicItemRepository
	RingRepository          repositories.RingRepository
	AmmoRepository          repositories.AmmoRepository
	SpellScrollRepository   repositories.SpellScrollRepository
	ContainerRepository     repositories.ContainerRepository
	TreasureRepository      repositories.TreasureRepository
	InventoryRepository     repositories.InventoryRepository
	ClassRepository         repositories.ClassRepository
	SpellCastingRepository  repositories.SpellCastingRepository
	WeaponMasteryRepository repositories.WeaponMasteryRepository

	ClassService       *services.ClassService
	EncumbranceService *services.EncumbranceService
	SpellService       *services.SpellService
	ACService          *services.ACService
	WeaponStatsService *services.WeaponStatsService

	UserController          *controllers.UserController
	CharacterController     *controllers.CharacterController
	SpellController         *controllers.SpellController
	ArmorController         *controllers.ArmorController
	WeaponController        *controllers.WeaponController
	EquipmentController     *controllers.EquipmentController
	ShieldController        *controllers.ShieldController
	PotionController        *controllers.PotionController
	MagicItemController     *controllers.MagicItemController
	RingController          *controllers.RingController
	AmmoController          *controllers.AmmoController
	SpellScrollController   *controllers.SpellScrollController
	ContainerController     *controllers.ContainerController
	AuthController          *controllers.AuthController
	TreasureController      *controllers.TreasureController
	InventoryController     *controllers.InventoryController
	SpellCastingController  *controllers.SpellCastingController
	ACController            *controllers.ACController
	WeaponMasteryController *controllers.WeaponMasteryController
	WeaponStatsController   *controllers.WeaponStatsController

	Templates      *template.Template
	SessionManager *scs.SessionManager
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
		"firstChar": func(s string) string {
			if len(s) > 0 {
				return string(s[0])
			}
			return ""
		},
	})

	// Parse templates
	tmpl, err = tmpl.ParseGlob(tmplPath)
	if err != nil {
		logger.Error("Failed to parse templates: %v", err)
		return nil, err
	}
	logger.Debug("Templates loaded successfully")

	// Set up session manager
	sessionManager := scs.New()
	logger.Debug("Setting up session store")
	sessionManager.Store = sqlite3store.New(db)
	logger.Debug("Session store initialized")
	sessionManager.Lifetime = 72 * time.Hour    // 3 days
	sessionManager.IdleTimeout = 24 * time.Hour // Reset after 1 day of inactivity
	sessionManager.Cookie.Name = "hyperborea_session"
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.Secure = false // Set to true in production with HTTPS
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode

	// Initialize repositories
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
	classRepo := repositories.NewSQLCClassRepository(db)
	spellCastingRepo := repositories.NewSQLCSpellCastingRepository(db)
	weaponMasteryRepo := repositories.NewSQLCWeaponMasteryRepository(db)

	// Initialize services
	classService := services.NewClassService(
		classRepo,
		inventoryRepo,
		armorRepo,
	)

	encumbranceService := services.NewEncumbranceService(
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
	)

	spellService := services.NewSpellService(
		spellRepo,
		spellCastingRepo,
		characterRepo,
		classRepo,
		classService,
		encumbranceService,
	)
	acService := services.NewACService(
		inventoryRepo,
		characterRepo,
		armorRepo,
		shieldRepo,
		encumbranceService,
	)
	weaponMasteryController := controllers.NewWeaponMasteryController(
		weaponMasteryRepo,
		characterRepo,
		weaponRepo,
	)

	weaponStatsService := services.NewWeaponStatsService(
		inventoryRepo,
		characterRepo,
		weaponRepo,
		weaponMasteryRepo,
	)

	classService.SetEncumbranceService(encumbranceService)

	// Initialize controllers with session manager
	authController := controllers.NewAuthController(userRepo, tmpl, sessionManager)
	userController := controllers.NewUserController(userRepo, tmpl)
	characterController := controllers.NewCharacterController(characterRepo, userRepo, classService, tmpl, sessionManager)
	spellController := controllers.NewSpellController(spellRepo, tmpl)
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
		encumbranceService,
		tmpl,
	)
	spellCastingController := controllers.NewSpellCastingController(spellService)
	acController := controllers.NewACController(acService)
	weaponStatsController := controllers.NewWeaponStatsController(weaponStatsService)
	logger.Info("Application initialized successfully")

	return &App{
		DB:                      db,
		UserRepository:          userRepo,
		CharacterRepository:     characterRepo,
		SpellRepository:         spellRepo,
		ArmorRepository:         armorRepo,
		WeaponRepository:        weaponRepo,
		EquipmentRepository:     equipmentRepo,
		ShieldRepository:        shieldRepo,
		PotionRepository:        potionRepo,
		MagicItemRepository:     magicItemRepo,
		RingRepository:          ringRepo,
		AmmoRepository:          ammoRepo,
		SpellScrollRepository:   spellScrollRepo,
		ContainerRepository:     containerRepo,
		TreasureRepository:      treasureRepo,
		InventoryRepository:     inventoryRepo,
		ClassRepository:         classRepo,
		SpellCastingRepository:  spellCastingRepo,
		WeaponMasteryRepository: weaponMasteryRepo,

		ClassService:       classService,
		EncumbranceService: encumbranceService,
		SpellService:       spellService,
		ACService:          acService,
		WeaponStatsService: weaponStatsService,

		UserController:          userController,
		CharacterController:     characterController,
		SpellController:         spellController,
		ArmorController:         armorController,
		WeaponController:        weaponController,
		EquipmentController:     equipmentController,
		ShieldController:        shieldController,
		PotionController:        potionController,
		MagicItemController:     magicItemController,
		RingController:          ringController,
		AmmoController:          ammoController,
		SpellScrollController:   spellScrollController,
		ContainerController:     containerController,
		AuthController:          authController,
		TreasureController:      treasureController,
		InventoryController:     inventoryController,
		SpellCastingController:  spellCastingController,
		ACController:            acController,
		WeaponMasteryController: weaponMasteryController,
		WeaponStatsController:   weaponStatsController,

		Templates:      tmpl,
		SessionManager: sessionManager,
	}, nil
}

func (a *App) SetupRoutes() http.Handler {
	logger.Debug("Setting up application routes")
	r := chi.NewRouter()

	// Load and save session data for all routes
	r.Use(a.SessionManager.LoadAndSave)

	// CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Static files handler
	r.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=3600")
		path := strings.TrimPrefix(r.URL.Path, "/static")
		if path != r.URL.Path {
			r.URL.Path = path
		}
		middleware.CustomFileServer(http.Dir("./web/static")).ServeHTTP(w, r)
	})

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Public routes (no auth required)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// If user is authenticated, redirect to dashboard
		if a.SessionManager.GetBool(r.Context(), "isAuthenticated") {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}
		// Otherwise, render the home page
		a.Templates.ExecuteTemplate(w, "home", nil)
	})

	// Authentication routes (no auth required)
	r.Route("/auth", func(r chi.Router) {
		r.Get("/login-page", a.AuthController.RenderLoginPage)
		r.Get("/register-page", a.AuthController.RenderRegisterPage)
		r.Post("/login", a.AuthController.Login)
		r.Post("/register", a.AuthController.Register)
		r.Get("/logout", a.AuthController.Logout)
	})

	// Protected routes (auth required)
	// Create authenticated router group
	authRouter := chi.NewRouter()
	authRouter.Use(a.requireAuthentication)

	authRouter.Get("/dashboard", a.CharacterController.RenderDashboard)

	// Protected web routes
	authRouter.Get("/settings", a.UserController.RenderSettingsPage)
	authRouter.Get("/characters/create", a.CharacterController.RenderCreateForm)
	authRouter.Get("/characters/view/{id}", a.CharacterController.RenderCharacterDetail)
	authRouter.Get("/characters/{id}/edit", a.CharacterController.RenderEditForm)

	// API routes requiring authentication
	authRouter.Route("/api", func(r chi.Router) {
		// User routes
		r.Route("/users", func(r chi.Router) {
			r.Get("/", a.UserController.ListUsers)
			r.Post("/", a.UserController.CreateUser)
			r.Get("/{id}", a.UserController.GetUser)
			r.Put("/{id}", a.UserController.UpdateUser)
			r.Delete("/{id}", a.UserController.DeleteUser)
			r.Get("/{id}/characters", a.CharacterController.GetCharactersByUser)
		})

		// Settings route
		r.Put("/user/settings", a.UserController.UpdateUserSettings)

		// Character routes
		r.Route("/characters", func(r chi.Router) {
			r.Get("/", a.CharacterController.ListCharacters)
			r.Post("/", a.CharacterController.CreateCharacter)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", a.CharacterController.GetCharacter)
				r.Put("/", a.CharacterController.UpdateCharacter)
				r.Delete("/", a.CharacterController.DeleteCharacter)
				r.Patch("/hp", a.CharacterController.UpdateCharacterHP)
				r.Post("/modify-hp", a.CharacterController.ModifyCharacterHP)
				r.Patch("/xp", a.CharacterController.UpdateCharacterXP)
				r.Get("/class-data", a.CharacterController.GetCharacterClassData)
				r.Get("/equipment-status", a.InventoryController.GetEquipmentStatus)
				r.Get("/combat-equipment", a.InventoryController.GetCombatEquipment)
				r.Get("/ac", a.ACController.GetCharacterAC)
				r.Get("/weapon-stats", a.WeaponStatsController.GetCharacterWeaponStats)

				r.Route("/weapon-masteries", func(r chi.Router) {
					r.Get("/", a.WeaponMasteryController.GetWeaponMasteriesByCharacter)
					r.Post("/", a.WeaponMasteryController.AddWeaponMastery)
					r.Put("/{weaponBaseName}", a.WeaponMasteryController.UpdateWeaponMastery)
					r.Delete("/{weaponBaseName}", a.WeaponMasteryController.DeleteWeaponMastery)
					r.Get("/available", a.WeaponMasteryController.GetAvailableWeaponsForMastery)
				})

				// Encumbrance routes
				r.Route("/encumbrance", func(r chi.Router) {
					r.Get("/", a.InventoryController.GetEncumbranceStatus)
					r.Post("/recalculate", a.InventoryController.RecalculateEncumbrance)
					r.Put("/capacity", a.InventoryController.UpdateInventoryCapacity)
				})

				// Spell routes
				r.Route("/spells", func(r chi.Router) {
					r.Get("/", a.SpellCastingController.GetCharacterSpellsInfo)
					r.Post("/known", a.SpellCastingController.AddKnownSpell)
					r.Delete("/known/{spellId}", a.SpellCastingController.RemoveKnownSpell)
					r.Post("/prepared", a.SpellCastingController.PrepareSpell)
					r.Delete("/prepared/{spellId}", a.SpellCastingController.UnprepareSpell)
					r.Delete("/prepared", a.SpellCastingController.ClearPreparedSpells)
					r.Post("/prepared/all", a.SpellCastingController.PrepareAllSpells)
					r.Get("/learnable", a.SpellCastingController.GetSpellsLearnableOnLevelUp)
					r.Post("/initial", a.SpellCastingController.AddInitialSpellsForNewCharacter)
				})
			})
		})

		// Game data routes
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

	// Mount authenticated router
	r.Mount("/", authRouter)

	// Apply recovery and logging middleware
	handler := middleware.RecoveryMiddleware(
		middleware.LoggingMiddleware(r),
	)

	logger.Info("Routes set up successfully")
	return handler
}

// Authentication middleware
func (a *App) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if user is authenticated
		userID := a.SessionManager.GetInt64(r.Context(), "userID")
		if userID == 0 {
			// Check if it's an API request or browser request
			if strings.HasPrefix(r.URL.Path, "/api/") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			} else {
				http.Redirect(w, r, "/auth/login-page", http.StatusSeeOther)
			}
			return
		}

		// Add isAuthenticated flag to the request context for templates
		ctx := r.Context()
		a.SessionManager.Put(ctx, "isAuthenticated", true)

		// User is authenticated, continue
		next.ServeHTTP(w, r)
	})
}

// Add context data for templates
func (a *App) addContextData(r *http.Request) map[string]interface{} {
	data := make(map[string]interface{})
	data["IsAuthenticated"] = a.SessionManager.GetBool(r.Context(), "isAuthenticated")

	if data["IsAuthenticated"].(bool) {
		userID := a.SessionManager.GetInt64(r.Context(), "userID")
		if userID > 0 {
			user, err := a.UserRepository.GetUser(r.Context(), userID)
			if err == nil {
				data["User"] = user
			}
		}
	}

	return data
}

func (a *App) Shutdown() {
	if a.DB != nil {
		logger.Info("Closing database connection...")
		a.DB.Close()
	}
}
