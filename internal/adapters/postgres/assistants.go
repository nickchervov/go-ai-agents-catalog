package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
)

func (r *PostgresRepository) Assistants(ctx context.Context, page, pageSize int, search string, includeInactive bool, categoryID *uuid.UUID) ([]domain.AssistantsList, int, error) {
	rows := make([]assistantsListRow, 0)
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
		a.updated_at
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
	if err := r.db.SelectContext(ctx, &rows, query, search, searchPattern, pageSize, offset, includeInactive, categoryID); err != nil {
		return nil, 0, fmt.Errorf("getting assistants from db: %w", err)
	}

	var total int
	totalQuery := `
	SELECT 
		COUNT(id) 
	FROM assistants a
	WHERE
		($4::uuid IS NULL OR a.category_id = $4::uuid)
			AND
		($3 = true OR a.is_activated = true)
			AND
		($1 = '' OR a.name ILIKE $2 OR a.description ILIKE $2)
	`
	if err := r.db.GetContext(ctx, &total, totalQuery, search, searchPattern, includeInactive, categoryID); err != nil {
		return nil, 0, fmt.Errorf("getting total assistants from db: %w", err)
	}

	assistants := make([]domain.AssistantsList, 0, len(rows))
	for _, row := range rows {
		assistants = append(assistants, row.toDomain())
	}

	return assistants, total, nil
}
