//go:build integration

package postgres

import (
	"github.com/google/uuid"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
)

func (s *AgentsRepositoryTestSuite) TestE2E() {
	user := &domain.User{
		ID:    uuid.New(),
		Email: "admin@example.com",
		Role:  "admin",
	}
	_, err := s.repo.UpsertUser(s.ctx, user)
	s.Require().NoError(err)

	category := &domain.Category{
		ID:   uuid.New(),
		Name: "Еда",
	}

	_, err = s.repo.CreateCategory(s.ctx, category)
	s.Require().NoError(err)

	assistant := &domain.Assistant{
		ID:           uuid.New(),
		Name:         "Повар",
		Description:  "что-то",
		CategoryID:   category.ID,
		Model:        "mock-model",
		SystemPrompt: "Создай домашний рецепт по ингредиентам",
		ExUserPrompt: "курица, лаваш, помидоры, огурцы",
		IsActivated:  true,
	}

	_, err = s.repo.CreateAssistant(s.ctx, assistant)
	s.Require().NoError(err)

	runs := &domain.AssistantRun{
		ID:          uuid.New(),
		AssistantID: assistant.ID,
		UserID:      user.ID,
		UserPrompt:  "первый запрос",
		Result:      "классный запрос",
		Status:      "success",
	}

	err = s.repo.RunAssistant(s.ctx, runs)
	s.Require().NoError(err)

	output, total, err := s.repo.UserRuns(s.ctx, user.ID, 1, 10, "")
	s.Require().NoError(err)
	s.Require().NotNil(output)

	s.Assert().Equal(1, total)
	s.Assert().Len(output, 1)

	for _, r := range output {
		s.Assert().Equal(runs.ID, r.ID)
		s.Assert().Equal(runs.AssistantID, r.AssistantID)
		s.Assert().Equal(runs.UserID, r.UserID)
		s.Assert().Equal(runs.UserPrompt, r.UserPrompt)
		s.Assert().Equal(runs.Result, r.Result)
		s.Assert().Equal(runs.Status, r.Status)
	}
}
