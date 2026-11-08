//go:build integration

package postgres

import (
	"github.com/google/uuid"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
)

func (s *AgentsRepositoryTestSuite) TestRunAssistantUpsert() {
	user := &domain.User{
		ID:    uuid.New(),
		Email: "admin@example.com",
		Role:  "admin",
	}

	category := &domain.Category{
		ID:   uuid.New(),
		Name: "Еда",
	}

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

	_, err := s.repo.UpsertUser(s.ctx, user)
	s.Require().NoError(err)
	_, err = s.repo.CreateCategory(s.ctx, category)
	s.Require().NoError(err)
	_, err = s.repo.CreateAssistant(s.ctx, assistant)
	s.Require().NoError(err)

	run := &domain.AssistantRun{
		ID:          uuid.New(),
		AssistantID: assistant.ID,
		UserID:      user.ID,
		UserPrompt:  "курица, лаваш, помидоры, яйца",
		Status:      domain.Pending,
	}

	err = s.repo.RunAssistant(s.ctx, run)
	s.Require().NoError(err)

	runs, total, err := s.repo.AdminRuns(s.ctx, nil, 1, 10, "")
	s.Require().NoError(err)
	s.Require().Equal(1, total)
	s.Require().Len(runs, 1)

	for _, r := range runs {
		s.Require().Equal(domain.Pending, r.Status)
	}

	run.Status = domain.Success
	run.Result = "work done"

	err = s.repo.RunAssistant(s.ctx, run)
	s.Require().NoError(err)

	runs, total, err = s.repo.AdminRuns(s.ctx, nil, 1, 10, "")
	s.Require().NoError(err)
	s.Require().Equal(1, total)
	s.Require().Len(runs, 1)

	for _, r := range runs {
		s.Require().Equal(domain.Success, r.Status)
		s.Require().Equal("work done", r.Result)
	}
}
