package app

import (
	"database/sql"
	"html/template"
	"mordezzanV4/internal/controllers"
	"mordezzanV4/internal/logger"
	"mordezzanV4/internal/middleware"
	"mordezzanV4/internal/repositories"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi"
	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	DB                  *sql.DB
	UserRepository      repositories.UserRepository
	CharacterRepository repositories.CharacterRepository
	SpellRepository     repositories.SpellRepository
	ArmorRepository     repositories.ArmorRepository
	WeaponRepository    repositories.WeaponRepository
	EquipmentRepository repositories.EquipmentRepository
	ShieldRepository    repositories.ShieldRepository
	PotionRepository    repositories.PotionRepository
	MagicItemRepository repositories.MagicItemRepository
	UserController      *controllers.UserController
	CharacterController *controllers.CharacterController
	SpellController     *controllers.SpellController
	ArmorController     *controllers.ArmorController
	WeaponController    *controllers.WeaponController
	EquipmentController *controllers.EquipmentController
	ShieldController    *controllers.ShieldController
	PotionController    *controllers.PotionController
	MagicItemController *controllers.MagicItemController

	Templates *template.Template
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
	tmpl, err := template.ParseGlob(tmplPath)
	if err != nil {
		logger.Error("Failed to parse templates: %v", err)
		return nil, err
	}
	logger.Debug("Templates loaded successfully")

	userRepo := repositories.NewSQLCUserRepository(db)
	characterRepo := repositories.NewSQLCCharacterRepository(db)
	spellRepo := repositories.NewSQLCSpellRepository(db)
	armorRepo := repositories.NewSQLCArmorRepository(db)
	weaponRepo := repositories.NewSQLCWeaponRepository(db)
	equipmentRepo := repositories.NewSQLCEquipmentRepository(db)
	shieldRepo := repositories.NewSQLCShieldRepository(db)
	potionRepo := repositories.NewSQLCPotionRepository(db)
	magicItemRepo := repositories.NewSQLCMagicItemRepository(db)

	userController := controllers.NewUserController(userRepo, tmpl)
	characterController := controllers.NewCharacterController(characterRepo, userRepo, tmpl)
	spellController := controllers.NewSpellController(spellRepo, characterRepo, tmpl)
	armorController := controllers.NewArmorController(armorRepo, tmpl)
	weaponController := controllers.NewWeaponController(weaponRepo, tmpl)
	equipmentController := controllers.NewEquipmentController(equipmentRepo, tmpl)
	shieldController := controllers.NewShieldController(shieldRepo, tmpl)
	potionController := controllers.NewPotionController(potionRepo, tmpl)
	magicItemController := controllers.NewMagicItemController(magicItemRepo, tmpl)

	logger.Info("Application initialized successfully")

	return &App{
		DB:                  db,
		UserRepository:      userRepo,
		CharacterRepository: characterRepo,
		SpellRepository:     spellRepo,
		ArmorRepository:     armorRepo,
		WeaponRepository:    weaponRepo,
		EquipmentRepository: equipmentRepo,
		ShieldRepository:    shieldRepo,
		PotionRepository:    potionRepo,
		MagicItemRepository: magicItemRepo,
		UserController:      userController,
		CharacterController: characterController,
		SpellController:     spellController,
		ArmorController:     armorController,
		WeaponController:    weaponController,
		EquipmentController: equipmentController,
		ShieldController:    shieldController,
		PotionController:    potionController,
		MagicItemController: magicItemController,
		Templates:           tmpl,
	}, nil
}

func (a *App) SetupRoutes() http.Handler {
	logger.Debug("Setting up application routes")
	r := chi.NewRouter()

	r.Get("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))).ServeHTTP)

	r.Route("/users", func(r chi.Router) {
		r.Get("/", a.UserController.ListUsers)
		r.Post("/", a.UserController.CreateUser)
		r.Get("/{id}", a.UserController.GetUser)
		r.Put("/{id}", a.UserController.UpdateUser)
		r.Delete("/{id}", a.UserController.DeleteUser)
		r.Get("/{id}/characters", a.CharacterController.GetCharactersByUser)
	})

	r.Route("/characters", func(r chi.Router) {
		r.Get("/", a.CharacterController.ListCharacters)
		r.Post("/", a.CharacterController.CreateCharacter)
		r.Get("/{id}", a.CharacterController.GetCharacter)
		r.Put("/{id}", a.CharacterController.UpdateCharacter)
		r.Delete("/{id}", a.CharacterController.DeleteCharacter)
		r.Get("/{id}/spells", a.SpellController.GetSpellsByCharacter)
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

	handler := middleware.RecoveryMiddleware(
		middleware.LoggingMiddleware(
			middleware.AuthMiddleware(r),
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
