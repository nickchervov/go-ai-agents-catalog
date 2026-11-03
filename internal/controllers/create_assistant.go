package controllers

import (
	"encoding/json"
	"mime"
	"net/http"

	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
	"github.com/nickchervov/go-ai-agents-catalog/internal/dto"
)

func (h *AgentsHandler) CreateAssistant(w http.ResponseWriter, r *http.Request) {
	mediatype, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if mediatype != "application/json" {
		JSON(w, http.StatusUnsupportedMediaType, map[string]ErrorMessage{"error": {Code: domain.ErrCodeInvalidRequest, Message: "unsupported media type"}})
		return
	}

	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	decoder.DisallowUnknownFields()

	var input dto.CreateAssistantInput
	if err := decoder.Decode(&input); err != nil {
		JSON(w, http.StatusBadRequest, map[string]ErrorMessage{"error": {Code: domain.ErrCodeInvalidRequest, Message: "decoding request body: " + err.Error()}})
		return
	}

	output, err := h.svc.CreateAssistant(r.Context(), input)
	if err != nil {
		JSONError(r.Context(), w, err)
		return
	}

	JSON(w, http.StatusCreated, output)
}
