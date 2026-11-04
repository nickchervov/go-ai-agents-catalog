package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
)

func (r *PostgresRepository) AssistantByID(ctx context.Context, id uuid.UUID) (*domain.AssistantById, error) {
	var row assistantByIdRow
	query := `
	SELECT 
		a.id,
		a.name,
		a.description,
		a.category_id,
		c.category,
		a.model,
		a.system_prompt,
		a.ex_user_prompt,
		a.is_activated,
		a.created_at,
		a.updated_at 
	FROM assistants a 
	LEFT JOIN categories c on a.category_id = c.id
	WHERE a.id = $1`
	if err := r.db.GetContext(ctx, &row, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrAssistantNotFound
		}
		return nil, fmt.Errorf("getting assistant by id from db: %w", err)
	}
	assistant := row.toDomain()
	return &assistant, nil
}
