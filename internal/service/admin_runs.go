package service

import (
	"context"
	"fmt"

	"github.com/nickchervov/go-ai-agents-catalog/internal/dto"
)

func (s *AgentService) AdminRuns(ctx context.Context, input dto.AdminRunsInput) (*dto.AdminRunsResponse, error) {
	runs, err := s.repo.AdminRuns(ctx, input.AssistantID, input.Page, input.PageSize, input.Status)
	if err != nil {
		return nil, fmt.Errorf("admin runs: %w", err)
	}

	output, total := dto.ToAdminRunsOutput(runs)

	pagination := dto.PaginationOutput{
		Page:     input.Page,
		PageSize: input.PageSize,
		Total:    total,
	}

	resp := dto.ToAdminRunsResponse(output, pagination)
	return resp, nil
}
