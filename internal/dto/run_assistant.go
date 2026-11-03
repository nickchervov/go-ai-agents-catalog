package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
)

type RunAssistantInput struct {
	UserPrompt  string `json:"userPrompt"`
	AssistantId uuid.UUID
	UserID      uuid.UUID
}

type RunAssistantOutput struct {
	ID            uuid.UUID        `json:"id"`
	AssistantID   uuid.UUID        `json:"assistantId"`
	AssistantName string           `json:"assistantName,omitempty"`
	CategoryID    uuid.UUID        `json:"categoryId,omitempty"`
	CategoryName  string           `json:"categoryName,omitempty"`
	UserID        uuid.UUID        `json:"userId"`
	Model         string           `json:"model"`
	UserPrompt    string           `json:"userPrompt"`
	Output        string           `json:"output,omitempty"`
	Status        domain.RunStatus `json:"status"`
	Error         string           `json:"error,omitempty"`
	CreatedAt     time.Time        `json:"createdAt,omitempty"`
}
