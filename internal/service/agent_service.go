package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
)

type AgentsRepository interface {
	Categories(ctx context.Context) ([]domain.Category, error)
	CreateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error)
	Assistants(ctx context.Context, page, pageSize int, search string, includeInactive bool, categoryID *uuid.UUID) ([]domain.AssistantsList, int, error)
	CreateAssistant(ctx context.Context, assistant *domain.Assistant) (*domain.AssistantById, error)
	AssistantByID(ctx context.Context, id uuid.UUID) (*domain.AssistantById, error)
	UpdateAssistantByID(ctx context.Context, id uuid.UUID, assistant *domain.Assistant) (*domain.AssistantById, error)
	UpsertUser(ctx context.Context, user *domain.User) (*domain.User, error)
	RunAssistant(ctx context.Context, runs *domain.AssistantRun) error
	UserRuns(ctx context.Context, userId uuid.UUID, page, pageSize int, status domain.RunStatus) ([]domain.AssistantRunsList, int, error)
	AdminRuns(ctx context.Context, assistantID *uuid.UUID, page, pageSize int, status domain.RunStatus) ([]domain.AssistantRunsList, int, error)
}

type TokenCreator interface {
	CreateJWTToken(userID, role, email string) (string, error)
}

type LLM interface {
	SendMessage(ctx context.Context, model, systemPrompt, userPrompt string) (string, error)
	Models(ctx context.Context) ([]string, error)
}

type AgentService struct {
	repo         AgentsRepository
	tokenCreator TokenCreator
	llm          LLM
}

func NewAgentService(repo AgentsRepository, creator TokenCreator, llm LLM) *AgentService {
	return &AgentService{repo: repo, tokenCreator: creator, llm: llm}
}
