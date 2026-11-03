package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
)

type UpdateAssistantInput struct {
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	CategoryID   uuid.UUID `json:"categoryId"`
	Model        string    `json:"model"`
	SystemPrompt string    `json:"systemPrompt"`
	ExUserPrompt string    `json:"exampleUserPrompt"`
	IsActivated  *bool     `json:"isActive"`
}

type UpdateAssistantsOutput struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	CategoryID   uuid.UUID `json:"categoryId"`
	CategoryName string    `json:"categoryName"`
	Model        string    `json:"model"`
	SystemPrompt string    `json:"systemPrompt"`
	ExUserPrompt string    `json:"exampleUserPrompt"`
	IsActivated  bool      `json:"isActive"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func FromUpdateAssistantDTO(input UpdateAssistantInput) *domain.Assistant {
	return &domain.Assistant{
		Name:         input.Name,
		Description:  input.Description,
		CategoryID:   input.CategoryID,
		Model:        input.Model,
		SystemPrompt: input.SystemPrompt,
		ExUserPrompt: input.ExUserPrompt,
		IsActivated:  input.IsActivated,
	}
}

func ToUpdateAssistantDTO(assistant *domain.AssistantsResponse) *UpdateAssistantsOutput {
	return &UpdateAssistantsOutput{
		ID:           assistant.ID,
		Name:         assistant.Name,
		Description:  assistant.Description,
		CategoryID:   assistant.CategoryID,
		CategoryName: assistant.CategoryName,
		Model:        assistant.Model,
		SystemPrompt: assistant.SystemPrompt,
		ExUserPrompt: assistant.ExUserPrompt,
		IsActivated:  assistant.IsActivated,
		CreatedAt:    assistant.CreatedAt,
		UpdatedAt:    assistant.UpdatedAt,
	}
}
