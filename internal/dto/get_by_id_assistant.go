package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
)

type ByIDAssistantsInput struct {
	ID uuid.UUID
}

type ByIDAssistantOutput struct {
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

func ToGetByIDAssistantDTO(userRole string, assistant *domain.AssistantByIdResponse) *ByIDAssistantOutput {
	output := ByIDAssistantOutput{
		ID:           assistant.ID,
		Name:         assistant.Name,
		Description:  assistant.Description,
		CategoryID:   assistant.CategoryID,
		CategoryName: assistant.CategoryName,
		Model:        assistant.Model,
		ExUserPrompt: assistant.ExUserPrompt,
		IsActivated:  assistant.IsActivated,
		CreatedAt:    assistant.CreatedAt,
		UpdatedAt:    assistant.UpdatedAt,
	}
	if userRole == "admin" {
		output.SystemPrompt = assistant.SystemPrompt
	}
	return &output
}
