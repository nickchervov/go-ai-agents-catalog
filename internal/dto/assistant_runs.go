package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
)

type AssistantRunsInput struct {
	AssistantID *uuid.UUID
	Status      domain.RunStatus
	Page        int
	PageSize    int
}

type AssistantRunsOutput struct {
	ID            uuid.UUID        `json:"id"`
	AssistantID   uuid.UUID        `json:"assistantId"`
	AssistantName string           `json:"assistantName"`
	CategoryID    uuid.UUID        `json:"categoryId"`
	CategoryName  string           `json:"categoryName"`
	UserID        uuid.UUID        `json:"userId"`
	Model         string           `json:"model"`
	UserPrompt    string           `json:"userPrompt"`
	Output        string           `json:"output"`
	Status        domain.RunStatus `json:"status"`
	Error         string           `json:"error"`
	CreatedAt     time.Time        `json:"createdAt"`
}

type AssistantRunsResponse struct {
	Runs       []AssistantRunsOutput `json:"runs"`
	Pagination PaginationOutput      `json:"pagination"`
}

func ToAssistantRunsOutput(runs []domain.AssistantRunsList, pagination PaginationOutput) *AssistantRunsResponse {
	r := make([]AssistantRunsOutput, 0)
	for _, v := range runs {
		dto := AssistantRunsOutput{
			ID:            v.ID,
			AssistantID:   v.AssistantID,
			AssistantName: v.AssistantName,
			CategoryID:    v.CategoryID,
			CategoryName:  v.CategoryName,
			UserID:        v.UserID,
			Model:         v.Model,
			UserPrompt:    v.UserPrompt,
			Output:        v.Result,
			Status:        v.Status,
			Error:         v.Error,
			CreatedAt:     v.CreatedAt,
		}
		r = append(r, dto)
	}

	return &AssistantRunsResponse{
		Runs:       r,
		Pagination: pagination,
	}
}
