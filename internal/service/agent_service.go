package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
)

type AgentsRepository interface {
	Categories(ctx context.Context) ([]*domain.Category, error)
	CreateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error)
	Assistants(ctx context.Context, page, pageSize int, search string, includeInactive bool, categoryID *uuid.UUID) ([]*domain.AssistantsResponse, error)
	CreateAssistant(ctx context.Context, assistant *domain.Assistant) (*domain.AssistantsResponse, error)
	AssistantByID(ctx context.Context, id uuid.UUID) (*domain.AssistantByIdResponse, error)
	UpdateAssistantByID(ctx context.Context, id uuid.UUID, assistant *domain.Assistant) (*domain.AssistantsResponse, error)
	UpsertUser(ctx context.Context, user *domain.User) (*domain.User, error)
	RunAssistant(ctx context.Context, runs *domain.AssistantRun) error
	UserRuns(ctx context.Context, userId uuid.UUID, page, pageSize int, status domain.RunStatus) ([]domain.UserRunsResponse, error)
	AdminRuns(ctx context.Context, assistantID *uuid.UUID, page, pageSize int, status domain.RunStatus) ([]domain.AdminRunsResponse, error)
}

type TokenCreator interface {
	CreateJWTToken(userID, role, email string) (string, error)
}

type LLMSender interface {
	SendMessage(ctx context.Context, model, systemPrompt, userPrompt string) (string, error)
}

type AgentService struct {
	repo         AgentsRepository
	tokenCreator TokenCreator
	llm          LLMSender
}

func NewAgentService(repo AgentsRepository, creator TokenCreator, llm LLMSender) *AgentService {
	return &AgentService{repo: repo, tokenCreator: creator, llm: llm}
}
