package postgres

import (
	"context"
	"fmt"

	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
)

func (r *PostgresRepository) Categories(ctx context.Context) ([]domain.Category, error) {
	rows := make([]categoryRow, 0)

	query := `
	SELECT 
		id, 
		category, 
		description, 
		created_at 
	FROM categories`

	if err := r.db.SelectContext(ctx, &rows, query); err != nil {
		return nil, fmt.Errorf("getting categories from db: %w", err)
	}

	categories := make([]domain.Category, 0, len(rows))
	for _, row := range rows {
		categories = append(categories, row.toDomain())
	}
	return categories, nil
}
