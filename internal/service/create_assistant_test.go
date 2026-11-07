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

func TestCreateAssistant(t *testing.T) {
	adminID, err := uuid.Parse(domain.RoleDatas["admin"].UserID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	categoryID := uuid.New()

	createdAt := time.Now()

	tests := []struct {
		name           string
		repoErr        error
		input          dto.CreateAssistantInput
		expectedOutput *dto.CreateAssistantsOutput
		checkErr       func(error) bool
	}{
		{
			name:    "good case",
			repoErr: nil,
			input: dto.CreateAssistantInput{
				Name:         "Повар",
				Description:  "Помогает создать ресторанный рецепт блюда",
				CategoryID:   categoryID,
				Model:        "mock-model-v2",
				SystemPrompt: "Составь рецепт по ингредиентам в несколько шагов",
				ExUserPrompt: "курица, сыр, лепешка, салат",
				IsActivated:  nil,
			},
			expectedOutput: &dto.CreateAssistantsOutput{
				Name:         "Повар",
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
			name:    "category not found",
			repoErr: nil,
			input: dto.CreateAssistantInput{
				Name:         "Повар",
				Description:  "Помогает создать ресторанный рецепт блюда",
				CategoryID:   uuid.New(),
				Model:        "mock-model-v2",
				SystemPrompt: "Составь рецепт по ингредиентам в несколько шагов",
				ExUserPrompt: "курица, сыр, лепешка, салат",
				IsActivated:  nil,
			},
			expectedOutput: nil,
			checkErr: func(err error) bool {
				return errors.Is(err, domain.ErrCategoryNotFound)
			},
		},
		{
			name:    "validation error",
			repoErr: nil,
			input: dto.CreateAssistantInput{
				Name:         "",
				Description:  "Помогает создать ресторанный рецепт блюда",
				CategoryID:   categoryID,
				Model:        "mock-model-v2",
				SystemPrompt: "Составь рецепт по ингредиентам в несколько шагов",
				ExUserPrompt: "курица, сыр, лепешка, салат",
				IsActivated:  nil,
			},
			expectedOutput: nil,
			checkErr: func(err error) bool {
				var targetErr *domain.ValidationErr
				return errors.As(err, &targetErr)
			},
		},
		{
			name:    "internal error",
			repoErr: ErrInternal,
			input: dto.CreateAssistantInput{
				Name:         "Повар",
				Description:  "Помогает создать ресторанный рецепт блюда",
				CategoryID:   categoryID,
				Model:        "mock-model-v2",
				SystemPrompt: "Составь рецепт по ингредиентам в несколько шагов",
				ExUserPrompt: "курица, сыр, лепешка, салат",
				IsActivated:  nil,
			},
			expectedOutput: nil,
			checkErr: func(err error) bool {
				return errors.Is(err, ErrInternal)
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
			}
			svc := service.NewAgentService(mockRepo, jwt, mockLLM)

			output, err := svc.CreateAssistant(t.Context(), tt.input)
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
