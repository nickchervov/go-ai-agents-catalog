package service_test

import (
	"errors"
	"testing"

	"github.com/nickchervov/go-ai-agents-catalog/internal/dto"
	"github.com/nickchervov/go-ai-agents-catalog/internal/service"
	"github.com/nickchervov/go-ai-agents-catalog/mocks"
)

func TestCreateCategory(t *testing.T) {
	tests := []struct {
		name           string
		repoError      error
		input          dto.CreateCategoryInput
		expectedOutput *dto.CreateCategoryOutput
		checkError     func(error) bool
	}{
		{
			name:      "good case",
			repoError: nil,
			input: dto.CreateCategoryInput{
				Name:        "Еда",
				Description: "что-то",
			},
			expectedOutput: &dto.CreateCategoryOutput{
				Name:        "Еда",
				Description: "что-то",
			},
			checkError: func(err error) bool {
				return err == nil
			},
		},
		{
			name:      "internal error",
			repoError: ErrInternal,
			input: dto.CreateCategoryInput{
				Name:        "Еда",
				Description: "что-то",
			},
			expectedOutput: nil,
			checkError: func(err error) bool {
				return errors.Is(err, ErrInternal)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewAgentService(&mocks.MockRepo{Err: tt.repoError}, &mocks.MockJWT{}, &mocks.MockLLM{})

			output, err := svc.CreateCategory(t.Context(), tt.input)
			if err != nil {
				if !tt.checkError(err) {
					t.Fatalf("unexpected error: %v", err)
				}
			} else {
				if output.Name != tt.expectedOutput.Name {
					t.Fatalf("got: %v want: %v", output.Name, tt.expectedOutput.Name)
				}
				if output.Description != tt.expectedOutput.Description {
					t.Fatalf("got: %v want: %v", output.Description, tt.expectedOutput.Description)
				}
			}
		})
	}
}
