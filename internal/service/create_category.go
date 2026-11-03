package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
	"github.com/nickchervov/go-ai-agents-catalog/internal/dto"
)

func (s *AgentService) CreateCategory(ctx context.Context, input dto.CreateCategoryInput) (*dto.CreateCategoryOutput, error) {
	category := dto.FromCreateCategoriesDTO(input)

	category.ID = uuid.New()

	if err := domain.ValidateCategory(category); err != nil {
		return nil, err
	}

	newCategory, err := s.repo.CreateCategory(ctx, category)
	if err != nil {
		return nil, fmt.Errorf("creating category: %w", err)
	}

	output := dto.ToCreateCategoryDTO(newCategory)

	return output, nil
}
