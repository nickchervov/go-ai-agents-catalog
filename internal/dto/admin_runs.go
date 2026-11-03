package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
)

type AdminRunsInput struct {
	AssistantID *uuid.UUID
	Status      domain.RunStatus
	Page        int
	PageSize    int
}

type AdminRunsOutput struct {
	ID            uuid.UUID `json:"id"`
	AssistantID   uuid.UUID `json:"assistantId"`
	AssistantName string    `json:"assistantName"`
	CategoryID    uuid.UUID `json:"categoryId"`
	CategoryName  string    `json:"categoryName"`
	UserID        uuid.UUID `json:"userId"`
	Model         string    `json:"model"`
	UserPrompt    string    `json:"userPrompt"`
	Output        string    `json:"output"`
	Status        string    `json:"status"`
	Error         string    `json:"error"`
	CreatedAt     time.Time `json:"createdAt"`
}

type AdminRunsResponse struct {
	Runs       []UserRunsOutput `json:"runs"`
	Pagination PaginationOutput `json:"pagination"`
}

func ToAdminRunsOutput(runs []domain.AdminRunsResponse) ([]UserRunsOutput, int) {
	total := 0
	r := make([]UserRunsOutput, 0)

	for _, v := range runs {
		if total == 0 {
			total = v.Total
		}
		dto := UserRunsOutput{
			ID:            v.ID,
			AssistantID:   v.AssistantID,
			AssistantName: v.AssistantName,
			CategoryID:    v.CategoryID,
			CategoryName:  v.CategoryName,
			UserID:        v.UserID,
			Model:         v.Model,
			UserPrompt:    v.UserPrompt,
			Output:        v.Output,
			Status:        v.Status,
			Error:         v.Error,
			CreatedAt:     v.CreatedAt,
		}
		r = append(r, dto)
	}

	return r, total
}

func ToAdminRunsResponse(output []UserRunsOutput, pagination PaginationOutput) *AdminRunsResponse {
	return &AdminRunsResponse{
		Runs:       output,
		Pagination: pagination,
	}
}
