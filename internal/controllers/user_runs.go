package controllers

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
	"github.com/nickchervov/go-ai-agents-catalog/internal/dto"
)

func (h *AgentsHandler) UserRuns(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	if status != "pending" && status != "success" && status != "failed" && status != "" {
		JSON(w, http.StatusBadRequest, map[string]ErrorMessage{"error": {Code: domain.ErrCodeInvalidRequest, Message: "incorrect status"}})
		return
	}
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	if page < 1 || pageSize < 1 {
		JSON(w, http.StatusBadRequest, map[string]ErrorMessage{"error": {Code: domain.ErrCodeInvalidRequest, Message: "incorrect page or page size values"}})
		return
	}

	userId, err := uuid.Parse(r.Context().Value("userId").(string))
	if err != nil {
		JSON(w, http.StatusBadRequest, map[string]ErrorMessage{"error": {Code: domain.ErrCodeInvalidRequest, Message: domain.ErrIncorrectUserID.Error()}})
		return
	}

	input := dto.UserRunsInput{
		Status:   domain.RunStatus(status),
		Page:     page,
		PageSize: pageSize,
		UserID:   userId,
	}

	output, err := h.svc.UserRuns(r.Context(), input)
	if err != nil {
		JSONError(r.Context(), w, err)
		return
	}

	JSON(w, http.StatusOK, output)
}
