package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
)

func (r *PostgresRepository) Assistants(ctx context.Context, page, pageSize int, search string, includeInactive bool, categoryID *uuid.UUID) ([]*domain.AssistantsResponse, error) {
	assistants := make([]*domain.AssistantsResponse, 0)
	offset := (page - 1) * pageSize
	query := `
	SELECT 
		a.id,
		a.name,
		a.description,
		a.category_id,
		c.category as category,
		a.model,
		a.system_prompt,
		a.ex_user_prompt,
		a.is_activated,
		a.created_at,
		a.updated_at,
		COUNT(*) OVER() as total_count
	FROM assistants a
	LEFT JOIN categories c ON a.category_id = c.id
	WHERE
		($6::uuid IS NULL OR a.category_id = $6::uuid)
			AND
		($5 = true OR a.is_activated = true)
			AND
		($1 = '' OR a.name ILIKE $2 OR a.description ILIKE $2)
	ORDER BY a.id
	LIMIT $3 OFFSET $4`
	var searchPattern string
	if search != "" {
		searchPattern = "%" + search + "%"
	}
	if err := r.db.SelectContext(ctx, &assistants, query, search, searchPattern, pageSize, offset, includeInactive, categoryID); err != nil {
		return nil, fmt.Errorf("getting assistants from db: %w", err)
	}

	return assistants, nil
}
