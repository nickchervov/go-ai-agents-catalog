package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
)

type GetAssistantsInput struct {
	CategoryId      *uuid.UUID
	Search          string
	IncludeInactive bool
	Page            int
	PageSize        int
}

type GetAssistantsOutput struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	CategoryID   uuid.UUID `json:"categoryId"`
	CategoryName string    `json:"categoryName"`
	Model        string    `json:"model"`
	SystemPrompt string    `json:"systemPrompt"`
	ExUserPrompt string    `json:"exampleUserPrompt"`
	IsActivated  bool      `json:"isActive"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type PaginationOutput struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}

type GetAssistantsResponse struct {
	Assistants []*GetAssistantsOutput `json:"assistants"`
	Pagination PaginationOutput       `json:"pagination"`
}

func ToGetAssistantsDTO(userRole string, assistants []*domain.AssistantsResponse) ([]*GetAssistantsOutput, int) {
	output := make([]*GetAssistantsOutput, 0)
	total := 0
	for _, v := range assistants {
		if total == 0 {
			total = v.TotalCount
		}
		dto := GetAssistantsOutput{
			ID:           v.ID,
			Name:         v.Name,
			Description:  v.Description,
			CategoryID:   v.CategoryID,
			CategoryName: v.CategoryName,
			Model:        v.Model,
			ExUserPrompt: v.ExUserPrompt,
			IsActivated:  v.IsActivated,
			CreatedAt:    v.CreatedAt,
			UpdatedAt:    v.UpdatedAt,
		}
		if userRole == "admin" {
			dto.SystemPrompt = v.SystemPrompt
		}
		output = append(output, &dto)
	}
	return output, total
}

func ToGetAssistantsResponse(output []*GetAssistantsOutput, pagination PaginationOutput) *GetAssistantsResponse {
	return &GetAssistantsResponse{
		Assistants: output,
		Pagination: pagination,
	}
}
