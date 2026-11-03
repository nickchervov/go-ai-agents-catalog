package postgres

import (
	"context"
	"fmt"

	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
)

func (r *PostgresRepository) Categories(ctx context.Context) ([]*domain.Category, error) {
	categories := make([]*domain.Category, 0)

	query := `
	SELECT 
		id, 
		category, 
		description, 
		created_at 
	FROM categories`

	if err := r.db.SelectContext(ctx, &categories, query); err != nil {
		return nil, fmt.Errorf("getting categories from db: %w", err)
	}

	return categories, nil
}
