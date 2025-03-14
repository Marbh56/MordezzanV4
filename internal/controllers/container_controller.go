package controllers

import (
	"encoding/json"
	"errors"
	"html/template"
	apperrors "mordezzanV4/internal/errors"
	"mordezzanV4/internal/models"
	"mordezzanV4/internal/repositories"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
)

type ContainerController struct {
	containerRepo repositories.ContainerRepository
	tmpl          *template.Template
}

func NewContainerController(containerRepo repositories.ContainerRepository, tmpl *template.Template) *ContainerController {
	return &ContainerController{
		containerRepo: containerRepo,
		tmpl:          tmpl,
	}
}

func (c *ContainerController) GetContainer(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid container ID format"))
		return
	}
	container, err := c.containerRepo.GetContainer(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}
	accept := r.Header.Get("Accept")
	if strings.Contains(accept, "application/json") {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(container); err != nil {
			apperrors.HandleError(w, apperrors.NewInternalError(err))
		}
		return
	}
	if err := c.tmpl.ExecuteTemplate(w, "container.html", container); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *ContainerController) ListContainers(w http.ResponseWriter, r *http.Request) {
	containers, err := c.containerRepo.ListContainers(r.Context())
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(containers); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *ContainerController) CreateContainer(w http.ResponseWriter, r *http.Request) {
	var input models.CreateContainerInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid request body format"))
		return
	}
	if err := input.Validate(); err != nil {
		var validationErr *models.ValidationError
		if errors.As(err, &validationErr) {
			validationErrors := map[string]string{
				validationErr.Field: validationErr.Message,
			}
			apperrors.HandleValidationErrors(w, validationErrors)
			return
		}
		apperrors.HandleError(w, err)
		return
	}
	existingContainer, err := c.containerRepo.GetContainerByName(r.Context(), input.Name)
	if err == nil && existingContainer != nil {
		validationErrors := map[string]string{
			"name": "Container with this name already exists",
		}
		apperrors.HandleValidationErrors(w, validationErrors)
		return
	}
	id, err := c.containerRepo.CreateContainer(r.Context(), &input)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}
	container, err := c.containerRepo.GetContainer(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(container); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *ContainerController) UpdateContainer(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid container ID format"))
		return
	}
	var input models.UpdateContainerInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid request body format"))
		return
	}
	if err := input.Validate(); err != nil {
		var validationErr *models.ValidationError
		if errors.As(err, &validationErr) {
			validationErrors := map[string]string{
				validationErr.Field: validationErr.Message,
			}
			apperrors.HandleValidationErrors(w, validationErrors)
			return
		}
		apperrors.HandleError(w, err)
		return
	}
	existingContainer, err := c.containerRepo.GetContainerByName(r.Context(), input.Name)
	if err == nil && existingContainer != nil && existingContainer.ID != id {
		validationErrors := map[string]string{
			"name": "Container with this name already exists",
		}
		apperrors.HandleValidationErrors(w, validationErrors)
		return
	}
	if err := c.containerRepo.UpdateContainer(r.Context(), id, &input); err != nil {
		apperrors.HandleError(w, err)
		return
	}
	updatedContainer, err := c.containerRepo.GetContainer(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedContainer); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *ContainerController) DeleteContainer(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid container ID format"))
		return
	}
	if err := c.containerRepo.DeleteContainer(r.Context(), id); err != nil {
		apperrors.HandleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
