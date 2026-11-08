//go:build integration

package postgres

import (
	"github.com/google/uuid"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
)

func (s *AgentsRepositoryTestSuite) TestCreateAssistant_CategoryNotFound() {
	assistant := &domain.Assistant{
		ID:           uuid.New(),
		Name:         "Повар",
		Description:  "что-то",
		CategoryID:   uuid.New(),
		Model:        "mock-model",
		SystemPrompt: "что-то",
		ExUserPrompt: "что-то",
		IsActivated:  true,
	}

	output, err := s.repo.CreateAssistant(s.ctx, assistant)
	s.Require().Nil(output)
	s.Require().ErrorIs(err, domain.ErrCategoryNotFound)
}
