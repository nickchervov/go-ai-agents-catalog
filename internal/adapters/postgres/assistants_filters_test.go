//go:build integration

package postgres

import (
	"github.com/google/uuid"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
)

func (s *AgentsRepositoryTestSuite) TestAssistantsFilters() {
	categories := []*domain.Category{
		{
			ID:   uuid.New(),
			Name: "Еда",
		},
		{
			ID:   uuid.New(),
			Name: "Дом",
		},
	}

	for _, c := range categories {
		_, err := s.repo.CreateCategory(s.ctx, c)
		s.Require().NoError(err)
	}

	assistants := []*domain.Assistant{
		{
			ID:           uuid.New(),
			Name:         "Повар",
			Description:  "что-то",
			CategoryID:   categories[0].ID,
			Model:        "mock-model",
			SystemPrompt: "Создай домашний рецепт по ингредиентам",
			ExUserPrompt: "курица, лаваш, помидоры, огурцы",
			IsActivated:  true,
		},
		{
			ID:           uuid.New(),
			Name:         "Повар 123",
			Description:  "что-то",
			CategoryID:   categories[0].ID,
			Model:        "mock-model",
			SystemPrompt: "Создай домашний рецепт по ингредиентам",
			ExUserPrompt: "курица, лаваш, помидоры, огурцы",
			IsActivated:  false,
		},
		{
			ID:           uuid.New(),
			Name:         "Муж на час",
			Description:  "что-то",
			CategoryID:   categories[1].ID,
			Model:        "mock-model",
			SystemPrompt: "Помоги решить проблему по дому",
			ExUserPrompt: "как вбить гвоздь",
			IsActivated:  true,
		},
		{
			ID:           uuid.New(),
			Name:         "Муж на час 123",
			Description:  "что-то",
			CategoryID:   categories[1].ID,
			Model:        "mock-model",
			SystemPrompt: "Помоги решить проблему по дому",
			ExUserPrompt: "как вбить гвоздь",
			IsActivated:  false,
		},
	}

	for _, a := range assistants {
		_, err := s.repo.CreateAssistant(s.ctx, a)
		s.Require().NoError(err)
	}

	output, total, err := s.repo.Assistants(s.ctx, 1, 10, "", true, &categories[0].ID)
	s.Require().NoError(err)
	s.Assert().Equal(2, total)
	s.Assert().NotNil(output)
	s.Assert().Len(output, 2)

	output, total, err = s.repo.Assistants(s.ctx, 1, 1, "", true, &categories[0].ID)
	s.Require().NoError(err)
	s.Assert().Equal(2, total)
	s.Assert().NotNil(output)
	s.Assert().Len(output, 1)

	output, total, err = s.repo.Assistants(s.ctx, 1, 10, "", true, nil)
	s.Require().NoError(err)
	s.Assert().Equal(4, total)
	s.Assert().NotNil(output)
	s.Assert().Len(output, 4)

	output, total, err = s.repo.Assistants(s.ctx, 1, 2, "", true, nil)
	s.Require().NoError(err)
	s.Assert().Equal(4, total)
	s.Assert().NotNil(output)
	s.Assert().Len(output, 2)

	output, total, err = s.repo.Assistants(s.ctx, 1, 10, "", false, nil)
	s.Require().NoError(err)
	s.Assert().Equal(2, total)
	s.Assert().NotNil(output)
	s.Assert().Len(output, 2)

	output, total, err = s.repo.Assistants(s.ctx, 1, 10, "", false, &categories[0].ID)
	s.Require().NoError(err)
	s.Assert().Equal(1, total)
	s.Assert().NotNil(output)
	s.Assert().Len(output, 1)
}
