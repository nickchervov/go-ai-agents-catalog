package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
	"github.com/nickchervov/go-ai-agents-catalog/internal/dto"
)

func (s *AgentService) CreateAssistant(ctx context.Context, input dto.CreateAssistantInput) (*dto.CreateAssistantsOutput, error) {
	assistant := dto.FromCreateAssistantDTO(input)

	assistant.ID = uuid.New()

	if err := domain.ValidationAssistant(assistant); err != nil {
		return nil, err
	}

	newAssistant, err := s.repo.CreateAssistant(ctx, assistant)
	if err != nil {
		return nil, fmt.Errorf("creating assistant: %w", err)
	}

	output := dto.ToCreateAssistantDTO(newAssistant)
	return output, nil
}
