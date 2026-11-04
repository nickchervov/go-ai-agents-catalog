package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
	"github.com/nickchervov/go-ai-agents-catalog/internal/dto"
)

type AgentsService interface {
	AdminRuns(ctx context.Context, input dto.AssistantRunsInput) (*dto.AssistantRunsResponse, error)
	AssistantByID(ctx context.Context, input dto.ByIDAssistantsInput) (*dto.ByIDAssistantOutput, error)
	Assistants(ctx context.Context, input dto.GetAssistantsInput) (*dto.GetAssistantsResponse, error)
	Categories(ctx context.Context) (*dto.GetCategoriesResponse, error)
	CreateAssistant(ctx context.Context, input dto.CreateAssistantInput) (*dto.CreateAssistantsOutput, error)
	CreateCategory(ctx context.Context, input dto.CreateCategoryInput) (*dto.CreateCategoryOutput, error)
	DummyLogin(ctx context.Context, input dto.DummyLoginInput) (*dto.DummyLoginOutput, error)
	RunAssistant(ctx context.Context, input dto.RunAssistantInput) (*dto.RunAssistantOutput, error)
	UpdateAssistant(ctx context.Context, id uuid.UUID, input dto.UpdateAssistantInput) (*dto.UpdateAssistantsOutput, error)
	UserRuns(ctx context.Context, input dto.UserAssistantRunsInput) (*dto.AssistantRunsResponse, error)
}

type ErrorMessage struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type TokenValidator interface {
	ValidateToken(token string) (*domain.Claims, error)
}

type AgentsHandler struct {
	svc AgentsService
}

func NewAgentsHandler(svc AgentsService) *AgentsHandler {
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
		return http.StatusBadRequest, domain.ErrCodeCategoryNotFound, domain.ErrCategoryNotFound.Error()
	case errors.Is(err, domain.ErrEmptyUserPrompt):
		return http.StatusBadRequest, domain.ErrCodeInvalidRequest, domain.ErrEmptyUserPrompt.Error()
	case errors.Is(err, domain.ErrLLMProvider):
		return http.StatusBadGateway, domain.ErrCodeLLMProvider, domain.ErrLLMProvider.Error()
	case errors.Is(err, domain.ErrEmptyIsActivated):
		return http.StatusBadRequest, domain.ErrCodeInvalidRequest, domain.ErrEmptyIsActivated.Error()
	default:
		return http.StatusInternalServerError, domain.ErrCodeInternal, "internal error"
	}
}

func (h *AgentsHandler) Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
