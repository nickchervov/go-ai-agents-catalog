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

func TestUserRuns(t *testing.T) {
	userID, err := uuid.Parse(domain.RoleDatas["user"].UserID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assistantID := uuid.New()
	categoryID := uuid.New()
	assistantRunID := uuid.New()

	tests := []struct {
		name           string
		repoError      error
		input          dto.UserAssistantRunsInput
		expectedOutput *dto.AssistantRunsResponse
		checkError     func(error) bool
	}{
		{
			name:      "good case",
			repoError: nil,
			input: dto.UserAssistantRunsInput{
				Page:     1,
				PageSize: 10,
				UserID:   userID,
			},
			expectedOutput: &dto.AssistantRunsResponse{
				Runs: []dto.AssistantRunsOutput{
					{
						ID:            assistantRunID,
						AssistantID:   assistantID,
						AssistantName: "Повар",
						CategoryID:    categoryID,
						CategoryName:  "Еда",
						UserID:        userID,
						Model:         "mock-model",
						UserPrompt:    "курица, сыр, огурцы, лаваш",
						Output:        "Good prompt take result back",
						Status:        "success",
					},
				},
				Pagination: dto.PaginationOutput{
					Page:     1,
					PageSize: 10,
					Total:    1,
				},
			},
			checkError: func(err error) bool {
				return err == nil
			},
		},
		{
			name:      "internal error",
			repoError: ErrInternal,
			input: dto.UserAssistantRunsInput{
				Page:     1,
				PageSize: 10,
				UserID:   userID,
			},
			expectedOutput: nil,
			checkError: func(err error) bool {
				return errors.Is(err, ErrInternal)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.MockRepo{
				Err: tt.repoError,
				User: &domain.User{
					ID:        userID,
					Email:     "user@example.com",
					Role:      "user",
					CreatedAt: time.Now(),
				},
				Category: &domain.Category{
					ID:          categoryID,
					Name:        "Еда",
					Description: "что-то про еду",
					CreatedAt:   time.Now(),
				},
				Assistant: &domain.Assistant{
					ID:           assistantID,
					Name:         "Повар",
					Description:  "Помогает составить рецепт из ингредиентов",
					CategoryID:   categoryID,
					Model:        "mock-model",
					SystemPrompt: "Составь рецепт по ингредиентам в несколько шагов",
					ExUserPrompt: "курица, сыр, лепешка, салат",
					IsActivated:  true,
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				},
				AssistantRun: &domain.AssistantRun{
					ID:          assistantRunID,
					AssistantID: assistantID,
					UserID:      userID,
					UserPrompt:  "курица, сыр, огурцы, лаваш",
					Result:      "Good prompt take result back",
					Status:      "success",
					CreatedAt:   time.Now(),
				},
			}

			svc := service.NewAgentService(mockRepo, &mocks.MockJWT{}, &mocks.MockLLM{Result: "Good prompt take result back"})

			output, err := svc.UserRuns(t.Context(), tt.input)
			if err != nil {
				if !tt.checkError(err) {
					t.Fatalf("unexpected error: %v", err)
				}
			} else {
				if len(output.Runs) != len(tt.expectedOutput.Runs) {
					t.Fatalf("got len: %d want: %d", len(output.Runs), len(tt.expectedOutput.Runs))
				}
				for i, r := range output.Runs {
					if r.ID != tt.expectedOutput.Runs[i].ID {
						t.Fatalf("got: %v want: %v", r.ID, tt.expectedOutput.Runs[i].ID)
					}
					if r.AssistantID != tt.expectedOutput.Runs[i].AssistantID {
						t.Fatalf("got: %v want: %v", r.AssistantID, tt.expectedOutput.Runs[i].AssistantID)
					}
					if r.AssistantName != tt.expectedOutput.Runs[i].AssistantName {
						t.Fatalf("got: %v want: %v", r.AssistantName, tt.expectedOutput.Runs[i].AssistantName)
					}
					if r.CategoryID != tt.expectedOutput.Runs[i].CategoryID {
						t.Fatalf("got: %v want: %v", r.CategoryID, tt.expectedOutput.Runs[i].CategoryID)
					}
					if r.CategoryName != tt.expectedOutput.Runs[i].CategoryName {
						t.Fatalf("got: %v want: %v", r.CategoryName, tt.expectedOutput.Runs[i].CategoryName)
					}
					if r.UserID != tt.expectedOutput.Runs[i].UserID {
						t.Fatalf("got: %v want: %v", r.UserID, tt.expectedOutput.Runs[i].UserID)
					}
					if r.Model != tt.expectedOutput.Runs[i].Model {
						t.Fatalf("got: %v want: %v", r.Model, tt.expectedOutput.Runs[i].Model)
					}
					if r.UserPrompt != tt.expectedOutput.Runs[i].UserPrompt {
						t.Fatalf("got: %v want: %v", r.UserPrompt, tt.expectedOutput.Runs[i].UserPrompt)
					}
					if r.Output != tt.expectedOutput.Runs[i].Output {
						t.Fatalf("got: %v want: %v", r.Output, tt.expectedOutput.Runs[i].Output)
					}
					if r.Status != tt.expectedOutput.Runs[i].Status {
						t.Fatalf("got: %v want: %v", r.Status, tt.expectedOutput.Runs[i].Status)
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
