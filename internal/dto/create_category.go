package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
)

type CreateCategoryInput struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type CreateCategoryOutput struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
}

func FromCreateCategoriesDTO(input CreateCategoryInput) *domain.Category {
	return &domain.Category{Name: input.Name, Description: input.Description}
}

func ToCreateCategoryDTO(category *domain.Category) *CreateCategoryOutput {
	return &CreateCategoryOutput{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
	}
}
