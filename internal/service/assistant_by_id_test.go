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

func TestAssistantByID(t *testing.T) {
	adminID, err := uuid.Parse(domain.RoleDatas["admin"].UserID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	categoryID := uuid.New()
	assistantID := uuid.New()

	tests := []struct {
		name           string
		repoErr        error
		input          dto.ByIDAssistantsInput
		expectedOutput *dto.ByIDAssistantOutput
		checkErr       func(error) bool
	}{
		{
			name:    "good case admin",
			repoErr: nil,
			input: dto.ByIDAssistantsInput{
				ID:       assistantID,
				UserRole: "admin",
			},
			expectedOutput: &dto.ByIDAssistantOutput{
				Name:         "Повар",
				SystemPrompt: "Составь рецепт по ингредиентам в несколько шагов",
			},
			checkErr: func(err error) bool {
				return err == nil
			},
		},
		{
			name:    "good case user",
			repoErr: nil,
			input: dto.ByIDAssistantsInput{
				ID:       assistantID,
				UserRole: "user",
			},
			expectedOutput: &dto.ByIDAssistantOutput{
				Name: "Повар",
			},
			checkErr: func(err error) bool {
				return err == nil
			},
		},
		{
			name:    "assistant not found",
			repoErr: nil,
			input: dto.ByIDAssistantsInput{
				ID:       uuid.New(),
				UserRole: "admin",
			},
			expectedOutput: nil,
			checkErr: func(err error) bool {
				return errors.Is(err, domain.ErrAssistantNotFound)
			},
		},
		{
			name:    "internal error",
			repoErr: ErrInternal,
			input: dto.ByIDAssistantsInput{
				ID:       uuid.New(),
				UserRole: "admin",
			},
			expectedOutput: nil,
			checkErr: func(err error) bool {
				return errors.Is(err, ErrInternal)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.MockRepo{
				Err: tt.repoErr,
				User: &domain.User{
					ID:        adminID,
					Email:     "admin@example.com",
					Role:      "admin",
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
			}

			mockLLM := &mocks.MockLLM{
				Result: "Nice prompt take result back",
			}

			jwt := auth.NewJWTProvider([]byte("qwerty"), 10*time.Minute)

			svc := service.NewAgentService(mockRepo, jwt, mockLLM)

			output, err := svc.AssistantByID(t.Context(), tt.input)
			if err != nil {
				if !tt.checkErr(err) {
					t.Fatalf("unexpected error: %v", err)
				}
			} else {
				if output.Name != tt.expectedOutput.Name {
					t.Fatalf("got name: %s want: %s", output.Name, tt.expectedOutput.Name)
				}
				if output.SystemPrompt != tt.expectedOutput.SystemPrompt {
					t.Fatalf("got sys prompt: %s want: %s", output.SystemPrompt, tt.expectedOutput.SystemPrompt)
				}
			}
		})
	}
}
