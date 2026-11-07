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

func TestRunAssistant(t *testing.T) {
	adminID, err := uuid.Parse(domain.RoleDatas["admin"].UserID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	categoryID := uuid.New()
	assistantID := uuid.New()

	tests := []struct {
		name              string
		repoErr           error
		llmErr            error
		inactiveAssistant bool
		input             dto.RunAssistantInput
		expectedOutput    *dto.RunAssistantOutput
		checkErr          func(error) bool
	}{
		{
			name:    "good case",
			repoErr: nil,
			llmErr:  nil,
			input: dto.RunAssistantInput{
				UserPrompt:  "курица, сыр, лаваш,капуста",
				AssistantId: assistantID,
				UserID:      adminID,
			},
			expectedOutput: &dto.RunAssistantOutput{
				UserPrompt:  "курица, сыр, лаваш,капуста",
				AssistantID: assistantID,
				UserID:      adminID,
				Output:      "Nice prompt take result back",
			},
			checkErr: func(err error) bool {
				return err == nil
			},
		},
		{
			name:    "not found assistant",
			repoErr: nil,
			llmErr:  nil,
			input: dto.RunAssistantInput{
				UserPrompt:  "курица, сыр, лаваш, капуста",
				AssistantId: uuid.New(),
				UserID:      adminID,
			},
			expectedOutput: nil,
			checkErr: func(err error) bool {
				return errors.Is(err, domain.ErrAssistantNotFound)
			},
		},
		{
			name:    "empty user prompt",
			repoErr: nil,
			llmErr:  nil,
			input: dto.RunAssistantInput{
				UserPrompt:  "",
				AssistantId: assistantID,
				UserID:      adminID,
			},
			expectedOutput: nil,
			checkErr: func(err error) bool {
				return errors.Is(err, domain.ErrEmptyUserPrompt)
			},
		},
		{
			name:    "only spaces user prompt",
			repoErr: nil,
			llmErr:  nil,
			input: dto.RunAssistantInput{
				UserPrompt:  "        ",
				AssistantId: assistantID,
				UserID:      adminID,
			},
			expectedOutput: nil,
			checkErr: func(err error) bool {
				return errors.Is(err, domain.ErrEmptyUserPrompt)
			},
		},
		{
			name:              "assistant inactive",
			repoErr:           nil,
			llmErr:            nil,
			inactiveAssistant: true,
			input: dto.RunAssistantInput{
				UserPrompt:  "курица, сыр, лаваш,капуста",
				AssistantId: assistantID,
				UserID:      adminID,
			},
			expectedOutput: nil,
			checkErr: func(err error) bool {
				return errors.Is(err, domain.ErrAssistantInactive)
			},
		},
		{
			name:    "LLM provider error",
			repoErr: nil,
			llmErr:  ErrLLMProvider,
			input: dto.RunAssistantInput{
				UserPrompt:  "курица, сыр, лаваш,капуста",
				AssistantId: assistantID,
				UserID:      adminID,
			},
			expectedOutput: nil,
			checkErr: func(err error) bool {
				return errors.Is(err, domain.ErrLLMProvider)
			},
		},
		{
			name:    "internal error",
			repoErr: ErrInternal,
			llmErr:  nil,
			input: dto.RunAssistantInput{
				UserPrompt:  "курица, сыр, лаваш,капуста",
				AssistantId: assistantID,
				UserID:      adminID,
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

			if tt.inactiveAssistant {
				mockRepo.Assistant.IsActivated = false
			}

			mockLLM := &mocks.MockLLM{
				Err:    tt.llmErr,
				Result: "Nice prompt take result back",
			}

			jwt := auth.NewJWTProvider([]byte("qwerty"), 10*time.Minute)
			svc := service.NewAgentService(mockRepo, jwt, mockLLM)

			output, err := svc.RunAssistant(t.Context(), tt.input)
			if err != nil {
				if !tt.checkErr(err) {
					t.Fatalf("unexpected error: %v", err)
				}
			} else {
				if output.UserPrompt != tt.expectedOutput.UserPrompt {
					t.Fatalf("got user prompt: %s want: %s", output.UserPrompt, tt.expectedOutput.UserPrompt)
				}
				if output.AssistantID != tt.expectedOutput.AssistantID {
					t.Fatalf("got assistant id: %v want: %v", output.AssistantID, tt.expectedOutput.AssistantID)
				}
				if output.UserID != tt.expectedOutput.UserID {
					t.Fatalf("got user id: %v want: %v", output.UserID, tt.expectedOutput.UserID)
				}
				if output.Output != tt.expectedOutput.Output {
					t.Fatalf("got output: %s want: %s", output.Output, tt.expectedOutput.Output)
				}
			}
		})
	}
}
