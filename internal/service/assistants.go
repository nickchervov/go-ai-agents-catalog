package service

import (
	"context"
	"fmt"

	"github.com/nickchervov/go-ai-agents-catalog/internal/dto"
)

func (s *AgentService) Assistants(ctx context.Context, input dto.GetAssistantsInput) (*dto.GetAssistantsResponse, error) {
	assistants, err := s.repo.Assistants(ctx, input.Page, input.PageSize, input.Search, input.IncludeInactive, input.CategoryId)
	if err != nil {
		return nil, fmt.Errorf("getting assistants: %w", err)
	}

	userRole, ok := ctx.Value("role").(string)
	if !ok {
		return nil, fmt.Errorf("getting user role from ctx")
	}
	output, total := dto.ToGetAssistantsDTO(userRole, assistants)

	response := dto.ToGetAssistantsResponse(output, dto.PaginationOutput{Page: input.Page, PageSize: input.PageSize, Total: total})
	return response, nil
}
