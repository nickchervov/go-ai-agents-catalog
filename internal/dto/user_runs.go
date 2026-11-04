package dto

import (
	"github.com/google/uuid"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
)

type UserAssistantRunsInput struct {
	AssistantID *uuid.UUID
	Status      domain.RunStatus
	Page        int
	PageSize    int
	UserID      uuid.UUID
}
