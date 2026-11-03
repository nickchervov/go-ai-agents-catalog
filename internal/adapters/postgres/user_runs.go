package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
)

func (r *PostgresRepository) UserRuns(ctx context.Context, userId uuid.UUID, page, pageSize int, status domain.RunStatus) ([]domain.UserRunsResponse, error) {
	offset := (page - 1) * pageSize

	runs := make([]domain.UserRunsResponse, 0)

	query := `
	SELECT
		ar.id,
		ar.assistant_id,
		a.name as assistant_name,
		a.category_id,
		c.category as category_name,
		ar.user_id,
		a.model,
		ar.user_prompt,
		ar.result as output,
		ar.status,
		ar.error,
		ar.created_at,
		COUNT(*) OVER() as total
	FROM assistant_run ar
	LEFT JOIN assistants a ON ar.assistant_id = a.id
	LEFT JOIN categories c ON a.category_id = c.id
	WHERE 
		ar.user_id = $1 AND ($2 = '' OR ar.status = $2)
	LIMIT $3 OFFSET $4
	`

	if err := r.db.SelectContext(ctx, &runs, query, userId, status, pageSize, offset); err != nil {
		return nil, fmt.Errorf("getting user runs: %w", err)
	}

	return runs, nil
}
