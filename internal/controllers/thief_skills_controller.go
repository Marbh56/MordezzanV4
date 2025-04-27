package controllers

import (
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	apperrors "mordezzanV4/internal/errors"
	"mordezzanV4/internal/logger"
	"mordezzanV4/internal/models"
	"mordezzanV4/internal/services"

	"github.com/go-chi/chi"
)

// ThiefSkillsRepository interface for the controller
type ThiefSkillsRepository interface {
	GetThiefSkillsByLevel(ctx context.Context, level int64) ([]*models.ThiefSkillWithChance, error)
	GetEffectiveThiefLevel(ctx context.Context, class string, level int64) (int64, error)
	ApplyAttributeBonus(skills []*models.ThiefSkillWithChance, attributes map[string]int) []*models.ThiefSkillWithChance
}

type ThiefSkillsServiceInterface interface {
	GetThiefSkillsForCharacter(ctx context.Context, charClass string, level int64, attributes map[string]int) ([]*models.ThiefSkillWithChance, error)
}

type CharacterRepository interface {
	GetCharacter(ctx context.Context, id int64) (*models.Character, error)
}

type ThiefSkillsController struct {
	thiefSkillsRepo    ThiefSkillsRepository
	characterRepo      CharacterRepository
	thiefSkillsService *services.ThiefSkillsService
	Templates          *template.Template
}

func NewThiefSkillsController(
	thiefSkillsRepo ThiefSkillsRepository,
	characterRepo CharacterRepository,
	thiefSkillsService *services.ThiefSkillsService,
	tmpl *template.Template,
) *ThiefSkillsController {
	return &ThiefSkillsController{
		thiefSkillsRepo:    thiefSkillsRepo,
		characterRepo:      characterRepo,
		thiefSkillsService: thiefSkillsService,
		Templates:          tmpl,
	}
}

func (c *ThiefSkillsController) RegisterRoutes(r chi.Router) {
	r.Route("/api/thief-skills", func(r chi.Router) {
		r.Get("/{level}", c.GetThiefSkillsByLevel)
	})

	r.Route("/api/characters/{id}/thief-skills", func(r chi.Router) {
		r.Get("/", c.GetThiefSkillsForCharacter)
	})
}

// GetThiefSkillsByLevel returns thief skills for a specific level
func (c *ThiefSkillsController) GetThiefSkillsByLevel(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Starting GetThiefSkillsByLevel function")
	levelStr := chi.URLParam(r, "level")
	level, err := strconv.ParseInt(levelStr, 10, 64)
	if err != nil {
		logger.Error("Failed to parse level: %v", err)
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid level format"))
		return
	}
	logger.Debug("Fetching thief skills for level: %d", level)

	skills, err := c.thiefSkillsRepo.GetThiefSkillsByLevel(r.Context(), level)
	if err != nil {
		logger.Error("Failed to get thief skills: %v", err)
		apperrors.HandleError(w, err)
		return
	}
	logger.Debug("Successfully fetched %d thief skills", len(skills))

	// Apply attribute bonuses if provided
	dexParam := r.URL.Query().Get("dexterity")
	intParam := r.URL.Query().Get("intelligence")
	wisParam := r.URL.Query().Get("wisdom")

	if dexParam != "" || intParam != "" || wisParam != "" {
		logger.Debug("Applying attribute bonuses to skills")
		// Build attributes map
		attributes := make(map[string]int)

		if dexParam != "" {
			if dex, err := strconv.Atoi(dexParam); err == nil {
				attributes["DX"] = dex
				logger.Debug("Applied DX attribute: %d", dex)
			}
		}

		if intParam != "" {
			if in, err := strconv.Atoi(intParam); err == nil {
				attributes["IN"] = in
				logger.Debug("Applied IN attribute: %d", in)
			}
		}

		if wisParam != "" {
			if wis, err := strconv.Atoi(wisParam); err == nil {
				attributes["WS"] = wis
				logger.Debug("Applied WS attribute: %d", wis)
			}
		}

		// Apply attribute bonuses
		skills = c.thiefSkillsRepo.ApplyAttributeBonus(skills, attributes)
	}

	logger.Debug("Sending response with %d skills", len(skills))
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(skills); err != nil {
		logger.Error("Failed to encode response: %v", err)
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
	logger.Debug("Completed GetThiefSkillsByLevel function")
}

func (c *ThiefSkillsController) GetThiefSkillsForCharacter(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Starting GetThiefSkillsForCharacter function")
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Error("Failed to parse character ID: %v", err)
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid character ID format"))
		return
	}
	logger.Debug("Fetching thief skills for character ID: %d", id)

	// Get character details
	character, err := c.characterRepo.GetCharacter(r.Context(), id)
	if err != nil {
		logger.Error("Failed to get character: %v", err)
		apperrors.HandleError(w, err)
		return
	}
	logger.Debug("Successfully fetched character: %s (Class: %s, Level: %d)", character.Name, character.Class, character.Level)

	// Output the literal class string for debugging
	logger.Debug("Character class string: '%s'", character.Class)

	// Build attributes map
	attributes := map[string]int{
		"DX": int(character.Dexterity),
		"IN": int(character.Intelligence),
		"WS": int(character.Wisdom),
	}
	logger.Debug("Character attributes - DX: %d, IN: %d, WS: %d", attributes["DX"], attributes["IN"], attributes["WS"])

	// Convert character.Level (int) to int64
	level := int64(character.Level)

	// Get thief skills for the character
	skills, err := c.thiefSkillsService.GetThiefSkillsForCharacter(
		r.Context(),
		character.Class,
		level,
		attributes,
	)
	if err != nil {
		logger.Error("Failed to get thief skills for character: %v", err)
		apperrors.HandleError(w, err)
		return
	}
	logger.Debug("Successfully fetched %d thief skills for character", len(skills))

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(skills); err != nil {
		logger.Error("Failed to encode response: %v", err)
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
	logger.Debug("Completed GetThiefSkillsForCharacter function")
}

// RenderThiefSkillsTab returns HTML for the thief skills tab
func (c *ThiefSkillsController) RenderThiefSkillsTab(w http.ResponseWriter, r *http.Request, charID int64) error {
	logger.Debug("Starting RenderThiefSkillsTab function for character ID: %d", charID)

	// Get character details
	character, err := c.characterRepo.GetCharacter(r.Context(), charID)
	if err != nil {
		logger.Error("Failed to get character: %v", err)
		return err
	}
	logger.Debug("Successfully fetched character: %s (Class: %s, Level: %d)", character.Name, character.Class, character.Level)

	// Build attributes map
	attributes := map[string]int{
		"DX": int(character.Dexterity),
		"IN": int(character.Intelligence),
		"WS": int(character.Wisdom),
	}
	logger.Debug("Character attributes - DX: %d, IN: %d, WS: %d", attributes["DX"], attributes["IN"], attributes["WS"])

	// Convert character.Level (int) to int64
	level := int64(character.Level)

	// Get thief skills for the character
	skills, err := c.thiefSkillsService.GetThiefSkillsForCharacter(
		r.Context(),
		character.Class,
		level,
		attributes,
	)
	if err != nil {
		logger.Error("Failed to get thief skills for character: %v", err)
		return err
	}
	logger.Debug("Successfully fetched %d thief skills for character", len(skills))

	// Convert to map for template
	thiefSkillsMap := make(map[string]string)
	for _, skill := range skills {
		thiefSkillsMap[skill.Name] = skill.SuccessChance
	}

	// Set data for template
	data := map[string]interface{}{
		"ThiefSkills": thiefSkillsMap,
	}
	logger.Debug("Rendering thief skills tab with %d skills", len(thiefSkillsMap))

	// Render template
	w.Header().Set("Content-Type", "text/html")
	if err := c.Templates.ExecuteTemplate(w, "thief_skills_tab", data); err != nil {
		logger.Error("Failed to execute template: %v", err)
		return apperrors.NewInternalError(err)
	}
	logger.Debug("Completed RenderThiefSkillsTab function")

	return nil
}
