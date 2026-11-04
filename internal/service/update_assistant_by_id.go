package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
	"github.com/nickchervov/go-ai-agents-catalog/internal/dto"
)

func (s *AgentService) UpdateAssistant(ctx context.Context, id uuid.UUID, input dto.UpdateAssistantInput) (*dto.UpdateAssistantsOutput, error) {
	if input.IsActivated == nil {
		return nil, domain.ErrEmptyIsActivated
	}

	assistant := dto.FromUpdateAssistantDTO(input)
	assistant.UpdatedAt = time.Now()

	if err := domain.ValidationAssistant(assistant); err != nil {
		return nil, err
	}
	newAssistant, err := s.repo.UpdateAssistantByID(ctx, id, assistant)
	if err != nil {
		return nil, fmt.Errorf("updating assistant: %w", err)
	}

	output := dto.ToUpdateAssistantDTO(newAssistant)
	return output, nil
}
