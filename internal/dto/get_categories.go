package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
)

type GetCategoriesOutput struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

type GetCategoriesResponse struct {
	Categories []GetCategoriesOutput `json:"categories"`
}

func ToGetCategoriesDTO(categories []domain.Category) *GetCategoriesResponse {
	output := make([]GetCategoriesOutput, 0)
	for _, v := range categories {
		dto := GetCategoriesOutput{ID: v.ID, Name: v.Name, Description: v.Description, CreatedAt: v.CreatedAt}
		output = append(output, dto)
	}
	return &GetCategoriesResponse{Categories: output}
}
