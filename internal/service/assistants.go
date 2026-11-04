package service

import (
	"context"
	"fmt"

	"github.com/nickchervov/go-ai-agents-catalog/internal/dto"
)

func (s *AgentService) Assistants(ctx context.Context, input dto.GetAssistantsInput) (*dto.GetAssistantsResponse, error) {
	assistants, total, err := s.repo.Assistants(ctx, input.Page, input.PageSize, input.Search, input.IncludeInactive, input.CategoryId)
	if err != nil {
		return nil, fmt.Errorf("getting assistants: %w", err)
	}

	pagination := dto.PaginationOutput{
		Page:     input.Page,
		PageSize: input.PageSize,
		Total:    total,
	}

	response := dto.ToGetAssistantsDTO(input.UserRole, assistants, pagination)

	return response, nil
}
