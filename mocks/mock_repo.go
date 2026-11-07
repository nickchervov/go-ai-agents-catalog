package mocks

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
)

type MockRepo struct {
	Category     *domain.Category
	Assistant    *domain.Assistant
	AssistantRun *domain.AssistantRun
	User         *domain.User
	Err          error
}

func (r *MockRepo) Categories(ctx context.Context) ([]domain.Category, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	categories := make([]domain.Category, 0)

	if r.Category == nil {
		return categories, nil
	}

	categories = append(categories, *r.Category)
	return categories, nil
}
func (r *MockRepo) CreateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	category.CreatedAt = time.Now()
	r.Category = category

	return category, nil
}
func (r *MockRepo) Assistants(ctx context.Context, page, pageSize int, search string, includeInactive bool, categoryID *uuid.UUID) ([]domain.AssistantsList, int, error) {
	if r.Err != nil {
		return nil, 0, r.Err
	}

	assistants := make([]domain.AssistantsList, 0)

	if r.Assistant == nil {
		return assistants, 0, nil
	}

	if !includeInactive && !r.Assistant.IsActivated {
		return assistants, 0, nil
	}

	if categoryID != nil {
		if r.Assistant.CategoryID != *categoryID {
			return assistants, 0, nil
		}
	}

	if strings.TrimSpace(search) != "" {
		if !strings.Contains(r.Assistant.Name, search) && !strings.Contains(r.Assistant.Description, search) {
			return assistants, 0, nil
		}
	}

	list := domain.AssistantsList{
		Assistant:    *r.Assistant,
		CategoryName: r.Category.Name,
	}

	assistants = append(assistants, list)
	return assistants, len(assistants), nil
}
func (r *MockRepo) CreateAssistant(ctx context.Context, assistant *domain.Assistant) (*domain.AssistantById, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	if assistant.CategoryID != r.Category.ID {
		return nil, domain.ErrCategoryNotFound
	}

	assistant.CreatedAt = time.Now()
	assistant.UpdatedAt = time.Now()

	r.Assistant = assistant

	assistantById := domain.AssistantById{
		Assistant:    *assistant,
		CategoryName: r.Category.Name,
	}
	return &assistantById, nil
}
func (r *MockRepo) AssistantByID(ctx context.Context, id uuid.UUID) (*domain.AssistantById, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	if r.Assistant == nil {
		return nil, domain.ErrAssistantNotFound
	}

	if r.Assistant.ID != id {
		return nil, domain.ErrAssistantNotFound
	}

	assistantById := domain.AssistantById{
		Assistant:    *r.Assistant,
		CategoryName: r.Category.Name,
	}

	return &assistantById, nil
}
func (r *MockRepo) UpdateAssistantByID(ctx context.Context, id uuid.UUID, assistant *domain.Assistant) (*domain.AssistantById, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	if r.Assistant.ID != id {
		return nil, domain.ErrAssistantNotFound
	}

	if assistant.CategoryID != r.Category.ID {
		return nil, domain.ErrCategoryNotFound
	}

	r.Assistant.Name = assistant.Name
	r.Assistant.Description = assistant.Description
	r.Assistant.CategoryID = assistant.CategoryID
	r.Assistant.Model = assistant.Model
	r.Assistant.SystemPrompt = assistant.SystemPrompt
	r.Assistant.ExUserPrompt = assistant.ExUserPrompt
	r.Assistant.IsActivated = assistant.IsActivated
	r.Assistant.UpdatedAt = time.Now()

	assistantById := domain.AssistantById{
		Assistant:    *r.Assistant,
		CategoryName: r.Category.Name,
	}
	return &assistantById, nil
}
func (r *MockRepo) UpsertUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	user.CreatedAt = time.Now()
	r.User = user
	return user, nil
}
func (r *MockRepo) RunAssistant(ctx context.Context, runs *domain.AssistantRun) error {
	if r.Err != nil {
		return r.Err
	}
	if r.Assistant.ID != runs.AssistantID {
		return domain.ErrAssistantNotFound
	}

	r.AssistantRun = runs
	return nil
}
func (r *MockRepo) UserRuns(ctx context.Context, userId uuid.UUID, page, pageSize int, status domain.RunStatus) ([]domain.AssistantRunsList, int, error) {
	if r.Err != nil {
		return nil, 0, r.Err
	}
	assistantsRunList := make([]domain.AssistantRunsList, 0)
	if userId != r.User.ID {
		return assistantsRunList, 0, nil
	}
	list := domain.AssistantRunsList{
		AssistantRun:  *r.AssistantRun,
		AssistantName: r.Assistant.Name,
		CategoryID:    r.Category.ID,
		CategoryName:  r.Category.Name,
		Model:         r.Assistant.Model,
	}
	assistantsRunList = append(assistantsRunList, list)
	return assistantsRunList, len(assistantsRunList), nil
}
func (r *MockRepo) AdminRuns(ctx context.Context, assistantID *uuid.UUID, page, pageSize int, status domain.RunStatus) ([]domain.AssistantRunsList, int, error) {
	if r.Err != nil {
		return nil, 0, r.Err
	}
	assistantsRunList := make([]domain.AssistantRunsList, 0)

	list := domain.AssistantRunsList{
		AssistantRun:  *r.AssistantRun,
		AssistantName: r.Assistant.Name,
		CategoryID:    r.Category.ID,
		CategoryName:  r.Category.Name,
		Model:         r.Assistant.Model,
	}
	assistantsRunList = append(assistantsRunList, list)
	return assistantsRunList, len(assistantsRunList), nil
}
