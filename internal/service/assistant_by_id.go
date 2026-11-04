package service

import (
	"context"
	"fmt"

	"github.com/nickchervov/go-ai-agents-catalog/internal/dto"
)

func (s *AgentService) AssistantByID(ctx context.Context, input dto.ByIDAssistantsInput) (*dto.ByIDAssistantOutput, error) {
	assistant, err := s.repo.AssistantByID(ctx, input.ID)
	if err != nil {
		return nil, fmt.Errorf("getting assistant by id: %w", err)
	}
	output := dto.ToGetByIDAssistantDTO(input.UserRole, assistant)

	return output, nil
}
