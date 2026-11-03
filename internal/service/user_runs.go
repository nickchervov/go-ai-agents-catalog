package service

import (
	"context"
	"fmt"

	"github.com/nickchervov/go-ai-agents-catalog/internal/dto"
)

func (s *AgentService) UserRuns(ctx context.Context, input dto.UserRunsInput) (*dto.UserRunsResponse, error) {
	runs, err := s.repo.UserRuns(ctx, input.UserID, input.Page, input.PageSize, input.Status)
	if err != nil {
		return nil, fmt.Errorf("user runs: %w", err)
	}

	output, total := dto.ToUserRunsOutput(runs)

	pagination := dto.PaginationOutput{
		Page:     input.Page,
		PageSize: input.PageSize,
		Total:    total,
	}

	resp := dto.ToUserRunsResponse(output, pagination)

	return resp, nil
}
