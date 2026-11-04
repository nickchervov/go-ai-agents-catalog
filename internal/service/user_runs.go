package service

import (
	"context"
	"fmt"

	"github.com/nickchervov/go-ai-agents-catalog/internal/dto"
)

func (s *AgentService) UserRuns(ctx context.Context, input dto.UserAssistantRunsInput) (*dto.AssistantRunsResponse, error) {
	runs, total, err := s.repo.UserRuns(ctx, input.UserID, input.Page, input.PageSize, input.Status)
	if err != nil {
		return nil, fmt.Errorf("user runs: %w", err)
	}

	pagination := dto.PaginationOutput{
		Page:     input.Page,
		PageSize: input.PageSize,
		Total:    total,
	}

	resp := dto.ToAssistantRunsOutput(runs, pagination)

	return resp, nil
}
