//go:build integration

package postgres

import (
	"github.com/google/uuid"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
)

func (s *AgentsRepositoryTestSuite) TestCategoryAssistantRoundTrip_Success() {
	categoryID := uuid.New()
	category := &domain.Category{
		ID:          categoryID,
		Name:        "Еда",
		Description: "что-то",
	}
	_, err := s.repo.CreateCategory(s.ctx, category)
	s.Require().NoError(err, "unexpected error: %v", err)

	assistantID := uuid.New()
	assistant := &domain.Assistant{
		ID:           assistantID,
		Name:         "Повар",
		Description:  "Созадет домашнее блюдо про ингредиентам",
		CategoryID:   categoryID,
		Model:        "mock-model",
		SystemPrompt: "Сделай домашний рецепт по переданным тебе ингредиентам",
		ExUserPrompt: "курица, лаваш, салат, сыр",
		IsActivated:  true,
	}

	_, err = s.repo.CreateAssistant(s.ctx, assistant)
	s.Require().NoError(err, "unexpected error: %v", err)

	assistentByID, err := s.repo.AssistantByID(s.ctx, assistantID)
	s.Require().NoError(err, "unexpected error: %v", err)

	s.Assert().Equal(category.Name, assistentByID.CategoryName, "got: %v want: %v", assistentByID.CategoryName, category.Name)
	s.Assert().Equal(assistant.ID, assistentByID.Assistant.ID, "got: %v want: %v", assistentByID.Assistant.ID, assistant.ID)
	s.Assert().Equal(assistant.Name, assistentByID.Assistant.Name, "got: %v want: %v", assistentByID.Assistant.Name, assistant.Name)
	s.Assert().Equal(assistant.IsActivated, assistentByID.Assistant.IsActivated, "got: %v want: %v", assistentByID.Assistant.IsActivated, assistant.IsActivated)
	s.Assert().Equal(assistant.Description, assistentByID.Assistant.Description, "got: %v want: %v", assistentByID.Assistant.Name, assistant.Description)

}
