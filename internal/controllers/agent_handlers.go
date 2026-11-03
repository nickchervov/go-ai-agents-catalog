package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"mime"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
	"github.com/nickchervov/go-ai-agents-catalog/internal/dto"
	"github.com/nickchervov/go-ai-agents-catalog/internal/service"
)

type ErrorMessage struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type TokenValidator interface {
	ValidateToken(token string) (*domain.Claims, error)
}

type AgentsHandler struct {
	svc *service.AgentService
}

func NewAgentsHandler(svc *service.AgentService) *AgentsHandler {
	return &AgentsHandler{svc: svc}
}

func JSON(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(data)
}

func JSONError(ctx context.Context, w http.ResponseWriter, err error) {
	reqId := middleware.GetReqID(ctx)

	code, msgCode, message := DecodeErrCode(err)
	if code == http.StatusInternalServerError {
		slog.Error("internal error", "request_id", reqId, "err", err)
	}
	JSON(w, code, map[string]ErrorMessage{"error": {Code: msgCode, Message: message}})
}

func DecodeErrCode(err error) (int, string, string) {
	var validationErr *domain.ValidationErr
	switch {
	case errors.Is(err, domain.ErrIncorrectCategoryName):
		return http.StatusBadRequest, domain.ErrCodeInvalidRequest, domain.ErrIncorrectCategoryName.Error()
	case errors.As(err, &validationErr):
		return http.StatusBadRequest, domain.ErrCodeInvalidRequest, validationErr.Error()
	case errors.Is(err, domain.ErrAssistantNotFound):
		return http.StatusNotFound, domain.ErrCodeAssistantNotFound, domain.ErrAssistantNotFound.Error()
	case errors.Is(err, domain.ErrIncorrectID):
		return http.StatusBadRequest, domain.ErrCodeInvalidRequest, domain.ErrIncorrectID.Error()
	case errors.Is(err, domain.ErrIncorrectUserID):
		return http.StatusBadRequest, domain.ErrCodeInvalidRequest, domain.ErrIncorrectUserID.Error()
	case errors.Is(err, domain.ErrIncorrectAssistantID):
		return http.StatusBadRequest, domain.ErrCodeInvalidRequest, domain.ErrIncorrectAssistantID.Error()
	case errors.Is(err, domain.ErrIncorrectRole):
		return http.StatusBadRequest, domain.ErrCodeInvalidRequest, domain.ErrIncorrectRole.Error()
	case errors.Is(err, domain.ErrAssistantInactive):
		return http.StatusConflict, domain.ErrCodeAssistantInactive, domain.ErrAssistantInactive.Error()
	case errors.Is(err, domain.ErrCategoryNotFound):
		return http.StatusConflict, domain.ErrCodeCategoryNotFound, domain.ErrCategoryNotFound.Error()
	case errors.Is(err, domain.ErrEmptyUserPrompt):
		return http.StatusBadRequest, domain.ErrCodeInvalidRequest, domain.ErrEmptyUserPrompt.Error()
	case errors.Is(err, domain.ErrLLMProvider):
		return http.StatusBadGateway, domain.ErrCodeLLMProvider, domain.ErrLLMProvider.Error()
	default:
		return http.StatusInternalServerError, domain.ErrCodeInternal, "internal error"
	}
}

func (h *AgentsHandler) DummyLogin(w http.ResponseWriter, r *http.Request) {
	mediatype, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if mediatype != "application/json" {
		JSON(w, http.StatusUnsupportedMediaType, map[string]ErrorMessage{"error": {Code: domain.ErrCodeInvalidRequest, Message: "unsupported media type"}})
		return
	}

	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	decoder.DisallowUnknownFields()

	var input dto.DummyLoginInput
	if err := decoder.Decode(&input); err != nil {
		JSON(w, http.StatusBadRequest, map[string]ErrorMessage{"error": {Code: domain.ErrCodeInvalidRequest, Message: "decoding request body: " + err.Error()}})
		return
	}

	output, err := h.svc.DummyLogin(r.Context(), input)
	if err != nil {
		JSONError(r.Context(), w, err)
		return
	}
	JSON(w, http.StatusOK, output)
}

func (h *AgentsHandler) Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
