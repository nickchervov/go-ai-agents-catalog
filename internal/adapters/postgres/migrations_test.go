//go:build integration

package postgres

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
)

func (s *AgentsRepositoryTestSuite) TestMigrations() {
	_, err := s.db.ExecContext(s.ctx, `
	INSERT INTO assistants (id, name, model, system_prompt, category_id)
		VALUES ($1,'n','m','sp', NULL )`, uuid.New())

	var pgErr *pgconn.PgError
	s.Require().ErrorAs(err, &pgErr)
	s.Assert().Equal("23502", pgErr.Code)

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

	_, err = s.repo.UpsertUser(s.ctx, user)
	s.Require().NoError(err)
	_, err = s.repo.CreateCategory(s.ctx, category)
	s.Require().NoError(err)
	_, err = s.repo.CreateAssistant(s.ctx, assistant)
	s.Require().NoError(err)

	_, err = s.db.ExecContext(s.ctx, `
	INSERT INTO assistant_run (id, assistant_id, user_id, user_prompt, status)
		VALUES ($1,$2,NULL,'sp', 'pending' )`, uuid.New(), assistant.ID)

	s.Require().ErrorAs(err, &pgErr)
	s.Assert().Equal("23502", pgErr.Code)
}
