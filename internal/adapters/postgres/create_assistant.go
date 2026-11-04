package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
)

func (r *PostgresRepository) CreateAssistant(ctx context.Context, assistant *domain.Assistant) (*domain.AssistantById, error) {
	var row assistantByIdRow
	query := `
	WITH inserted_assistant AS (
		INSERT INTO assistants 
			(id, name, description, category_id, model, system_prompt, ex_user_prompt, is_activated)
		VALUES
			($1, $2, $3, $4, $5, $6, $7,$8)
		RETURNING 
			id, name, description, category_id, model,system_prompt, ex_user_prompt, is_activated, created_at, updated_at
	)
			
	SELECT
		ia.id,
		ia.name,
		ia.description,
		ia.category_id,
		c.category as category,
		ia.model,
		ia.system_prompt,
		ia.ex_user_prompt,
		ia.is_activated, 
		ia.created_at, 
		ia.updated_at
	FROM inserted_assistant as ia
	LEFT JOIN categories AS c ON ia.category_id = c.id`

	if err := r.db.GetContext(ctx, &row, query, assistant.ID, assistant.Name, assistant.Description,
		assistant.CategoryID, assistant.Model, assistant.SystemPrompt, assistant.ExUserPrompt, assistant.IsActivated); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return nil, domain.ErrCategoryNotFound
		}
		return nil, fmt.Errorf("inserting assistant into db: %w", err)
	}

	newAssistant := row.toDomain()
	return &newAssistant, nil
}
