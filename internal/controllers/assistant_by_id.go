package controllers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
	"github.com/nickchervov/go-ai-agents-catalog/internal/dto"
)

func (h *AgentsHandler) AssistantByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "assistantId"))
	if err != nil {
		JSON(w, http.StatusBadRequest, map[string]ErrorMessage{"error": {Code: domain.ErrCodeInvalidRequest, Message: domain.ErrIncorrectID.Error()}})
		return
	}
	userRole, _ := r.Context().Value("role").(string)

	input := dto.ByIDAssistantsInput{
		ID:       id,
		UserRole: userRole,
	}

	output, err := h.svc.AssistantByID(r.Context(), input)
	if err != nil {
		JSONError(r.Context(), w, err)
		return
	}

	JSON(w, http.StatusOK, output)
}
