package controllers

import (
	"encoding/json"
	"mime"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
	"github.com/nickchervov/go-ai-agents-catalog/internal/dto"
)

func (h *AgentsHandler) RunAssistant(w http.ResponseWriter, r *http.Request) {
	mediatype, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if mediatype != "application/json" {
		JSON(w, http.StatusUnsupportedMediaType, map[string]ErrorMessage{"error": {Code: domain.ErrCodeInvalidRequest, Message: "unsupported media type"}})
		return
	}

	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	decoder.DisallowUnknownFields()

	var input dto.RunAssistantInput
	if err := decoder.Decode(&input); err != nil {
		JSON(w, http.StatusBadRequest, map[string]ErrorMessage{"error": {Code: domain.ErrCodeInvalidRequest, Message: "decoding request body: " + err.Error()}})
		return
	}

	userID, err := uuid.Parse(r.Context().Value("userId").(string))
	if err != nil {
		JSON(w, http.StatusBadRequest, map[string]ErrorMessage{"error": {Code: domain.ErrCodeInvalidRequest, Message: "parsing user id"}})
		return
	}
	assistantID, err := uuid.Parse(chi.URLParam(r, "assistantId"))
	if err != nil {
		JSON(w, http.StatusBadRequest, map[string]ErrorMessage{"error": {Code: domain.ErrCodeInvalidRequest, Message: "parsing assistant id"}})
		return
	}

	input.UserID = userID
	input.AssistantId = assistantID

	output, err := h.svc.RunAssistant(r.Context(), input)
	if err != nil {
		JSONError(r.Context(), w, err)
		return
	}

	JSON(w, http.StatusCreated, output)
}
