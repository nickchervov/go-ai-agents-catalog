package service

import (
	"context"
	"fmt"

	"github.com/nickchervov/go-ai-agents-catalog/internal/dto"
)

func (s *AgentService) AdminRuns(ctx context.Context, input dto.AssistantRunsInput) (*dto.AssistantRunsResponse, error) {
	runs, total, err := s.repo.AdminRuns(ctx, input.AssistantID, input.Page, input.PageSize, input.Status)
	if err != nil {
		return nil, fmt.Errorf("admin runs: %w", err)
	}

	pagination := dto.PaginationOutput{
		Page:     input.Page,
		PageSize: input.PageSize,
		Total:    total,
	}

	resp := dto.ToAssistantRunsOutput(runs, pagination)
	return resp, nil
}
