package service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nickchervov/go-ai-agents-catalog/internal/auth"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
	"github.com/nickchervov/go-ai-agents-catalog/internal/dto"
	"github.com/nickchervov/go-ai-agents-catalog/internal/service"
	"github.com/nickchervov/go-ai-agents-catalog/mocks"
)

func TestUpdateAssistant(t *testing.T) {
	adminID, err := uuid.Parse(domain.RoleDatas["admin"].UserID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	categoryID := uuid.New()
	assistantID := uuid.New()
	createdAt := time.Now()

	trueVal := true

	tests := []struct {
		name           string
		repoErr        error
		id             uuid.UUID
		input          dto.UpdateAssistantInput
		expectedOutput *dto.UpdateAssistantsOutput
		checkErr       func(error) bool
	}{
		{
			name:    "good case",
			repoErr: nil,
			id:      assistantID,
			input: dto.UpdateAssistantInput{
				Name:         "Повар Мишлен",
				Description:  "Помогает создать ресторанный рецепт блюда",
				CategoryID:   categoryID,
				Model:        "mock-model-v2",
				SystemPrompt: "Составь рецепт по ингредиентам в несколько шагов",
				ExUserPrompt: "курица, сыр, лепешка, салат",
				IsActivated:  &trueVal,
			},
			expectedOutput: &dto.UpdateAssistantsOutput{
				Name:         "Повар Мишлен",
				Description:  "Помогает создать ресторанный рецепт блюда",
				CategoryID:   categoryID,
				CategoryName: "Еда",
				Model:        "mock-model-v2",
				SystemPrompt: "Составь рецепт по ингредиентам в несколько шагов",
				ExUserPrompt: "курица, сыр, лепешка, салат",
				IsActivated:  true,
			},
			checkErr: func(err error) bool {
				return err == nil
			},
		},
		{
			name:    "empty is activated",
			repoErr: nil,
			id:      assistantID,
			input: dto.UpdateAssistantInput{
				Name:         "Повар Мишлен",
				Description:  "Помогает создать ресторанный рецепт блюда",
				CategoryID:   categoryID,
				Model:        "mock-model-v2",
				SystemPrompt: "Составь рецепт по ингредиентам в несколько шагов",
				ExUserPrompt: "курица, сыр, лепешка, салат",
				IsActivated:  nil,
			},
			expectedOutput: nil,
			checkErr: func(err error) bool {
				return errors.Is(err, domain.ErrEmptyIsActivated)
			},
		},
		{
			name:    "assistant not found",
			repoErr: nil,
			id:      uuid.New(),
			input: dto.UpdateAssistantInput{
				Name:         "Повар Мишлен",
				Description:  "Помогает создать ресторанный рецепт блюда",
				CategoryID:   categoryID,
				Model:        "mock-model-v2",
				SystemPrompt: "Составь рецепт по ингредиентам в несколько шагов",
				ExUserPrompt: "курица, сыр, лепешка, салат",
				IsActivated:  &trueVal,
			},
			expectedOutput: nil,
			checkErr: func(err error) bool {
				return errors.Is(err, domain.ErrAssistantNotFound)
			},
		},
		{
			name:    "category not found",
			repoErr: nil,
			id:      assistantID,
			input: dto.UpdateAssistantInput{
				Name:         "Повар Мишлен",
				Description:  "Помогает создать ресторанный рецепт блюда",
				CategoryID:   uuid.New(),
				Model:        "mock-model-v2",
				SystemPrompt: "Составь рецепт по ингредиентам в несколько шагов",
				ExUserPrompt: "курица, сыр, лепешка, салат",
				IsActivated:  &trueVal,
			},
			expectedOutput: nil,
			checkErr: func(err error) bool {
				return errors.Is(err, domain.ErrCategoryNotFound)
			},
		},
		{
			name:    "internal error",
			repoErr: ErrInternal,
			id:      assistantID,
			input: dto.UpdateAssistantInput{
				Name:         "Повар Мишлен",
				Description:  "Помогает создать ресторанный рецепт блюда",
				CategoryID:   categoryID,
				Model:        "mock-model-v2",
				SystemPrompt: "Составь рецепт по ингредиентам в несколько шагов",
				ExUserPrompt: "курица, сыр, лепешка, салат",
				IsActivated:  &trueVal,
			},
			expectedOutput: nil,
			checkErr: func(err error) bool {
				return errors.Is(err, ErrInternal)
			},
		},
		{
			name:    "validation error",
			repoErr: nil,
			id:      assistantID,
			input: dto.UpdateAssistantInput{
				Name:         "",
				Description:  "Помогает создать ресторанный рецепт блюда",
				CategoryID:   categoryID,
				Model:        "mock-model-v2",
				SystemPrompt: "Составь рецепт по ингредиентам в несколько шагов",
				ExUserPrompt: "курица, сыр, лепешка, салат",
				IsActivated:  &trueVal,
			},
			expectedOutput: nil,
			checkErr: func(err error) bool {
				var targetErr *domain.ValidationErr
				return errors.As(err, &targetErr)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLLM := &mocks.MockLLM{
				Result: "Nice prompt take result back",
			}

			jwt := auth.NewJWTProvider([]byte("qwerty"), 10*time.Minute)

			mockRepo := &mocks.MockRepo{
				Err: tt.repoErr,
				User: &domain.User{
					ID:        adminID,
					Email:     "admin@example.com",
					Role:      "admin",
					CreatedAt: createdAt,
				},
				Category: &domain.Category{
					ID:          categoryID,
					Name:        "Еда",
					Description: "что-то про еду",
					CreatedAt:   createdAt,
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
					CreatedAt:    createdAt,
					UpdatedAt:    createdAt,
				},
			}

			svc := service.NewAgentService(mockRepo, jwt, mockLLM)

			output, err := svc.UpdateAssistant(t.Context(), tt.id, tt.input)
			if err != nil {
				if !tt.checkErr(err) {
					t.Fatalf("unexpected error: %v", err)
				}
			} else {
				if output.ID == uuid.Nil {
					t.Fatal("id is empty")
				}
				if output.Name != tt.expectedOutput.Name {
					t.Fatalf("got: %s want: %s output: %v", output.Name, tt.expectedOutput.Name, output)
				}
				if output.Description != tt.expectedOutput.Description {
					t.Fatalf("got: %s want: %s output: %v", output.Description, tt.expectedOutput.Description, output)
				}
				if output.CategoryID != tt.expectedOutput.CategoryID {
					t.Fatalf("got: %v want: %v output: %v", output.CategoryID, tt.expectedOutput.CategoryID, output)
				}
				if output.CategoryName != tt.expectedOutput.CategoryName {
					t.Fatalf("got: %v want: %v output: %v", output.CategoryName, tt.expectedOutput.CategoryName, output)
				}
				if output.Model != tt.expectedOutput.Model {
					t.Fatalf("got: %s want: %s output: %v", output.Model, tt.expectedOutput.Model, output)
				}
				if output.SystemPrompt != tt.expectedOutput.SystemPrompt {
					t.Fatalf("got: %s want: %s output: %v", output.SystemPrompt, tt.expectedOutput.SystemPrompt, output)
				}
				if output.ExUserPrompt != tt.expectedOutput.ExUserPrompt {
					t.Fatalf("got: %s want: %s output: %v", output.ExUserPrompt, tt.expectedOutput.ExUserPrompt, output)
				}
				if output.IsActivated != tt.expectedOutput.IsActivated {
					t.Fatalf("got: %t want: %t output: %v", output.IsActivated, tt.expectedOutput.IsActivated, output)
				}
			}
		})
	}
}
