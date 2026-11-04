package controllers

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
	"github.com/nickchervov/go-ai-agents-catalog/internal/dto"
)

func (h *AgentsHandler) Assistants(w http.ResponseWriter, r *http.Request) {
	page := 1
	strPage := r.URL.Query().Get("page")
	if strPage != "" {
		intPage, err := strconv.Atoi(strPage)
		if err != nil {
			JSON(w, http.StatusBadRequest, map[string]ErrorMessage{"error": {Code: domain.ErrCodeInvalidRequest, Message: "incorrect page"}})
			return
		}
		page = intPage
	}
	pageSize := 10
	strPageSize := r.URL.Query().Get("pageSize")
	if strPageSize != "" {
		intPageSize, err := strconv.Atoi(strPageSize)
		if err != nil {
			JSON(w, http.StatusBadRequest, map[string]ErrorMessage{"error": {Code: domain.ErrCodeInvalidRequest, Message: "incorrect page size"}})
			return
		}
		pageSize = intPageSize
	}
	if pageSize > 100 {
		pageSize = 100
	}
	if page < 1 || pageSize < 1 {
		JSON(w, http.StatusBadRequest, map[string]ErrorMessage{"error": {Code: domain.ErrCodeInvalidRequest, Message: "incorrect page or page size values"}})
		return
	}

	search := r.URL.Query().Get("q")

	userRole, _ := r.Context().Value("role").(string)

	input := dto.GetAssistantsInput{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
		UserRole: userRole,
	}

	includeInactive := r.URL.Query().Get("includeInactive")
	if includeInactive != "" && includeInactive == "true" {
		input.IncludeInactive = true
	}

	if role, _ := r.Context().Value("role").(string); role != "admin" {
		input.IncludeInactive = false
	}

	categoryId := r.URL.Query().Get("categoryId")
	if categoryId != "" {
		parseUUID, err := uuid.Parse(categoryId)
		if err != nil {
			JSON(w, http.StatusBadRequest, map[string]ErrorMessage{"error": {Code: domain.ErrCodeInvalidRequest, Message: "parsing category id"}})
			return
		}
		input.CategoryId = &parseUUID
	}

	output, err := h.svc.Assistants(r.Context(), input)
	if err != nil {
		JSONError(r.Context(), w, err)
		return
	}

	JSON(w, http.StatusOK, output)
}
