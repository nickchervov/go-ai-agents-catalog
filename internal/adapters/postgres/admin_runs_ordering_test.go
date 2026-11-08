//go:build integration

package postgres

import (
	"time"

	"github.com/google/uuid"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
)

func (s *AgentsRepositoryTestSuite) TestAdminRunsOrdering() {
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

	runs := []*domain.AssistantRun{
		{
			ID:          uuid.New(),
			AssistantID: assistant.ID,
			UserID:      user.ID,
			UserPrompt:  "первый запрос",
			Result:      "классный запрос",
			Status:      "success",
		},
		{
			ID:          uuid.New(),
			AssistantID: assistant.ID,
			UserID:      user.ID,
			UserPrompt:  "второй запрос",
			Result:      "классный запрос",
			Status:      "success",
		},
		{
			ID:          uuid.New(),
			AssistantID: assistant.ID,
			UserID:      user.ID,
			UserPrompt:  "третий запрос",
			Result:      "классный запрос",
			Status:      "success",
		},
	}

	_, err := s.repo.UpsertUser(s.ctx, user)
	s.Require().NoError(err)
	_, err = s.repo.CreateCategory(s.ctx, category)
	s.Require().NoError(err)
	_, err = s.repo.CreateAssistant(s.ctx, assistant)
	s.Require().NoError(err)

	for _, r := range runs {
		err := s.repo.RunAssistant(s.ctx, r)
		time.Sleep(1 * time.Second)

		s.Require().NoError(err)
	}

	output, total, err := s.repo.AdminRuns(s.ctx, nil, 1, 10, domain.Success)
	s.Require().NoError(err)
	s.Require().NotNil(output)

	s.Assert().Equal(3, total)

	for i := 1; i < len(output); i++ {
		s.Assert().Less(output[i].CreatedAt, output[i-1].CreatedAt)
	}

	output, total, err = s.repo.UserRuns(s.ctx, uuid.New(), 1, 10, "")
	s.Require().NoError(err)
	s.Require().NotNil(output)

	s.Assert().Equal(0, total)
	s.Assert().Empty(output)
}
