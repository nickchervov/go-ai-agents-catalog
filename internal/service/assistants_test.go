package service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
	"github.com/nickchervov/go-ai-agents-catalog/internal/dto"
	"github.com/nickchervov/go-ai-agents-catalog/internal/service"
	"github.com/nickchervov/go-ai-agents-catalog/mocks"
)

func TestAssistants(t *testing.T) {
	assistantID := uuid.New()
	categoryID := uuid.New()

	tests := []struct {
		name           string
		repoError      error
		input          dto.GetAssistantsInput
		expectedOutput *dto.GetAssistantsResponse
		checkError     func(error) bool
		isActivated    bool
	}{
		{
			name:      "good case admin",
			repoError: nil,
			input: dto.GetAssistantsInput{
				CategoryId:      nil,
				Search:          "",
				IncludeInactive: true,
				Page:            1,
				PageSize:        10,
				UserRole:        "admin",
			},
			expectedOutput: &dto.GetAssistantsResponse{
				Assistants: []dto.GetAssistantsOutput{{
					ID:           assistantID,
					Name:         "Повар",
					Description:  "что-то",
					CategoryID:   categoryID,
					CategoryName: "Еда",
					Model:        "mock-model",
					SystemPrompt: "Сделай домашний рецепт по ингредиентам",
					ExUserPrompt: "курица, рис, макароны, сыр",
					IsActivated:  true,
				}},
				Pagination: dto.PaginationOutput{
					Page:     1,
					PageSize: 10,
					Total:    1,
				},
			},
			checkError: func(err error) bool {
				return err == nil
			},
			isActivated: true,
		},
		{
			name:      "good case user",
			repoError: nil,
			input: dto.GetAssistantsInput{
				CategoryId: nil,
				Search:     "",
				Page:       1,
				PageSize:   10,
				UserRole:   "user",
			},
			expectedOutput: &dto.GetAssistantsResponse{
				Assistants: []dto.GetAssistantsOutput{{
					ID:           assistantID,
					Name:         "Повар",
					Description:  "что-то",
					CategoryID:   categoryID,
					CategoryName: "Еда",
					Model:        "mock-model",
					ExUserPrompt: "курица, рис, макароны, сыр",
					IsActivated:  true,
				}},
				Pagination: dto.PaginationOutput{
					Page:     1,
					PageSize: 10,
					Total:    1,
				},
			},
			checkError: func(err error) bool {
				return err == nil
			},
			isActivated: true,
		},
		{
			name:      "include inactive",
			repoError: nil,
			input: dto.GetAssistantsInput{
				CategoryId:      nil,
				Search:          "",
				Page:            1,
				PageSize:        10,
				IncludeInactive: true,
				UserRole:        "admin",
			},
			expectedOutput: &dto.GetAssistantsResponse{
				Assistants: []dto.GetAssistantsOutput{{
					ID:           assistantID,
					Name:         "Повар",
					Description:  "что-то",
					CategoryID:   categoryID,
					CategoryName: "Еда",
					Model:        "mock-model",
					SystemPrompt: "Сделай домашний рецепт по ингредиентам",
					ExUserPrompt: "курица, рис, макароны, сыр",
					IsActivated:  false,
				}},
				Pagination: dto.PaginationOutput{
					Page:     1,
					PageSize: 10,
					Total:    1,
				},
			},
			checkError: func(err error) bool {
				return err == nil
			},
			isActivated: false,
		},
		{
			name:      "include search",
			repoError: nil,
			input: dto.GetAssistantsInput{
				CategoryId:      nil,
				Search:          "что",
				Page:            1,
				PageSize:        10,
				IncludeInactive: true,
				UserRole:        "admin",
			},
			expectedOutput: &dto.GetAssistantsResponse{
				Assistants: []dto.GetAssistantsOutput{{
					ID:           assistantID,
					Name:         "Повар",
					Description:  "что-то",
					CategoryID:   categoryID,
					CategoryName: "Еда",
					Model:        "mock-model",
					SystemPrompt: "Сделай домашний рецепт по ингредиентам",
					ExUserPrompt: "курица, рис, макароны, сыр",
					IsActivated:  false,
				}},
				Pagination: dto.PaginationOutput{
					Page:     1,
					PageSize: 10,
					Total:    1,
				},
			},
			checkError: func(err error) bool {
				return err == nil
			},
			isActivated: false,
		},
		{
			name:      "include category id",
			repoError: nil,
			input: dto.GetAssistantsInput{
				CategoryId:      &categoryID,
				Search:          "",
				Page:            1,
				PageSize:        10,
				IncludeInactive: true,
				UserRole:        "admin",
			},
			expectedOutput: &dto.GetAssistantsResponse{
				Assistants: []dto.GetAssistantsOutput{{
					ID:           assistantID,
					Name:         "Повар",
					Description:  "что-то",
					CategoryID:   categoryID,
					CategoryName: "Еда",
					Model:        "mock-model",
					SystemPrompt: "Сделай домашний рецепт по ингредиентам",
					ExUserPrompt: "курица, рис, макароны, сыр",
					IsActivated:  false,
				}},
				Pagination: dto.PaginationOutput{
					Page:     1,
					PageSize: 10,
					Total:    1,
				},
			},
			checkError: func(err error) bool {
				return err == nil
			},
			isActivated: false,
		},
		{
			name:      "internal error",
			repoError: ErrInternal,
			input: dto.GetAssistantsInput{
				CategoryId:      nil,
				Search:          "",
				Page:            1,
				PageSize:        10,
				IncludeInactive: true,
				UserRole:        "admin",
			},
			expectedOutput: nil,
			checkError: func(err error) bool {
				return errors.Is(err, ErrInternal)
			},
			isActivated: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.MockRepo{
				Err: tt.repoError,
				Category: &domain.Category{
					ID:          categoryID,
					Name:        "Еда",
					Description: "что-то",
					CreatedAt:   time.Now(),
				},
				Assistant: &domain.Assistant{
					ID:           assistantID,
					Name:         "Повар",
					Description:  "что-то",
					CategoryID:   categoryID,
					Model:        "mock-model",
					SystemPrompt: "Сделай домашний рецепт по ингредиентам",
					ExUserPrompt: "курица, рис, макароны, сыр",
					IsActivated:  tt.isActivated,
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				},
			}

			mockLLM := &mocks.MockLLM{}

			mockJWT := &mocks.MockJWT{}

			svc := service.NewAgentService(mockRepo, mockJWT, mockLLM)

			output, err := svc.Assistants(t.Context(), tt.input)
			if err != nil {
				if !tt.checkError(err) {
					t.Fatalf("unexpected error")
				}
			} else {
				if len(output.Assistants) != len(tt.expectedOutput.Assistants) {
					t.Fatalf("got len: %d want: %d", len(output.Assistants), len(tt.expectedOutput.Assistants))
				}
				for i := range output.Assistants {
					if output.Assistants[i].ID != tt.expectedOutput.Assistants[i].ID {
						t.Fatalf("got: %v want: %v", output.Assistants[i].ID, tt.expectedOutput.Assistants[i].ID)
					}
					if output.Assistants[i].Name != tt.expectedOutput.Assistants[i].Name {
						t.Fatalf("got: %v want: %v", output.Assistants[i].Name, tt.expectedOutput.Assistants[i].Name)
					}
					if output.Assistants[i].Description != tt.expectedOutput.Assistants[i].Description {
						t.Fatalf("got: %v want: %v", output.Assistants[i].Description, tt.expectedOutput.Assistants[i].Description)
					}
					if output.Assistants[i].CategoryID != tt.expectedOutput.Assistants[i].CategoryID {
						t.Fatalf("got: %v want: %v", output.Assistants[i].CategoryID, tt.expectedOutput.Assistants[i].CategoryID)
					}
					if output.Assistants[i].CategoryName != tt.expectedOutput.Assistants[i].CategoryName {
						t.Fatalf("got: %v want: %v", output.Assistants[i].CategoryName, tt.expectedOutput.Assistants[i].CategoryName)
					}
					if output.Assistants[i].Model != tt.expectedOutput.Assistants[i].Model {
						t.Fatalf("got: %v want: %v", output.Assistants[i].Model, tt.expectedOutput.Assistants[i].Model)
					}
					if output.Assistants[i].SystemPrompt != tt.expectedOutput.Assistants[i].SystemPrompt {
						t.Fatalf("got: %v want: %v", output.Assistants[i].SystemPrompt, tt.expectedOutput.Assistants[i].SystemPrompt)
					}
					if output.Assistants[i].ExUserPrompt != tt.expectedOutput.Assistants[i].ExUserPrompt {
						t.Fatalf("got: %v want: %v", output.Assistants[i].ExUserPrompt, tt.expectedOutput.Assistants[i].ExUserPrompt)
					}
					if output.Assistants[i].IsActivated != tt.expectedOutput.Assistants[i].IsActivated {
						t.Fatalf("got: %v want: %v", output.Assistants[i].IsActivated, tt.expectedOutput.Assistants[i].IsActivated)
					}
				}
				if output.Pagination.Page != tt.input.Page {
					t.Fatalf("got: %v want: %v", output.Pagination.Page, tt.input.Page)
				}

				if output.Pagination.PageSize != tt.input.PageSize {
					t.Fatalf("got: %v want: %v", output.Pagination.PageSize, tt.input.PageSize)
				}

				if output.Pagination.Total != 1 {
					t.Fatalf("got: %v want: 1", output.Pagination.Total)
				}
			}
		})
	}
}
