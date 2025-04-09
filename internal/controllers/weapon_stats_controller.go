package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	apperrors "mordezzanV4/internal/errors"
	"mordezzanV4/internal/services"

	"github.com/go-chi/chi"
)

type WeaponStatsController struct {
	weaponStatsService *services.WeaponStatsService
}

func NewWeaponStatsController(weaponStatsService *services.WeaponStatsService) *WeaponStatsController {
	return &WeaponStatsController{
		weaponStatsService: weaponStatsService,
	}
}

func (c *WeaponStatsController) GetCharacterWeaponStats(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid character ID format"))
		return
	}

	weaponStats, err := c.weaponStatsService.CalculateCharacterWeaponStats(r.Context(), id)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			apperrors.HandleError(w, apperrors.NewNotFound("character", id))
		} else {
			apperrors.HandleError(w, apperrors.NewInternalError(err))
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"weapon_stats": weaponStats,
	})
}
