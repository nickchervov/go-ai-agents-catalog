package service

import (
	"context"
	"fmt"

	"github.com/nickchervov/go-ai-agents-catalog/internal/dto"
)

func (s *AgentService) Categories(ctx context.Context) (*dto.GetCategoriesResponse, error) {
	categories, err := s.repo.Categories(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting categories: %w", err)
	}

	output := dto.ToGetCategoriesDTO(categories)

	return output, nil
}
