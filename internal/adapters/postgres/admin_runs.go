package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
)

func (r *PostgresRepository) AdminRuns(ctx context.Context, assistantID *uuid.UUID, page, pageSize int, status domain.RunStatus) ([]domain.AssistantRunsList, int, error) {
	offset := (page - 1) * pageSize
	rows := make([]assistantRunsListRow, 0)

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
		ar.result,
		ar.status,
		ar.error,
		ar.created_at
	FROM assistant_run ar
	LEFT JOIN assistants a ON ar.assistant_id = a.id
	LEFT JOIN categories c ON a.category_id = c.id
	WHERE 
		($1::uuid IS NULL OR ar.assistant_id = $1) AND ($2 = '' OR ar.status = $2)
	ORDER BY ar.created_at DESC, ar.id DESC
	LIMIT $3 OFFSET $4
	`

	if err := r.db.SelectContext(ctx, &rows, query, assistantID, status, pageSize, offset); err != nil {
		return nil, 0, fmt.Errorf("getting all runs from db: %w", err)
	}

	var total int
	totalQuery := `
	SELECT
		COUNT(id)
	FROM assistant_run ar
	WHERE
		($1::uuid IS NULL OR ar.assistant_id = $1) AND ($2 = '' OR ar.status = $2)`

	if err := r.db.GetContext(ctx, &total, totalQuery, assistantID, status); err != nil {
		return nil, 0, fmt.Errorf("getting total count of assistant runs: %w", err)
	}

	runs := make([]domain.AssistantRunsList, 0, len(rows))
	for _, row := range rows {
		runs = append(runs, row.toDomain())
	}

	return runs, total, nil
}
