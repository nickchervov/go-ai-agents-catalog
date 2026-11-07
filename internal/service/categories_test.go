package service_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
	"github.com/nickchervov/go-ai-agents-catalog/internal/dto"
	"github.com/nickchervov/go-ai-agents-catalog/internal/service"
	"github.com/nickchervov/go-ai-agents-catalog/mocks"
)

func TestCategories(t *testing.T) {
	categoryID := uuid.New()

	tests := []struct {
		name           string
		repoError      error
		expectedOutput *dto.GetCategoriesResponse
		checkError     func(error) bool
	}{
		{
			name:      "good case",
			repoError: nil,
			expectedOutput: &dto.GetCategoriesResponse{
				Categories: []dto.GetCategoriesOutput{
					{
						ID:          categoryID,
						Name:        "Еда",
						Description: "что-то",
					},
				},
			},
			checkError: func(err error) bool {
				return err == nil
			},
		},
		{
			name:           "internal error",
			repoError:      ErrInternal,
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
				Category: &domain.Category{
					ID:          categoryID,
					Name:        "Еда",
					Description: "что-то",
				},
			}

			svc := service.NewAgentService(mockRepo, &mocks.MockJWT{}, &mocks.MockLLM{})

			output, err := svc.Categories(t.Context())
			if err != nil {
				if !tt.checkError(err) {
					t.Fatalf("unexpected error: %v", err)
				}
			} else {
				if len(output.Categories) != len(tt.expectedOutput.Categories) {
					t.Fatalf("got len: %d want: %d", len(output.Categories), len(tt.expectedOutput.Categories))
				}

				for i := range output.Categories {
					if output.Categories[i].ID != tt.expectedOutput.Categories[i].ID {
						t.Fatalf("got: %v want: %v", output.Categories[i].ID, tt.expectedOutput.Categories[i].ID)
					}
					if output.Categories[i].Name != tt.expectedOutput.Categories[i].Name {
						t.Fatalf("got: %v want: %v", output.Categories[i].Name, tt.expectedOutput.Categories[i].Name)
					}
					if output.Categories[i].Description != tt.expectedOutput.Categories[i].Description {
						t.Fatalf("got: %v want: %v", output.Categories[i].Description, tt.expectedOutput.Categories[i].Description)
					}
				}
			}
		})
	}
}
