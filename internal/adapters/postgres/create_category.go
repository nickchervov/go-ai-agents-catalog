package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
)

func (r *PostgresRepository) CreateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	query := `
	INSERT INTO categories 
		(id,category, description) 
	VALUES 
		($1,$2,$3)
	ON CONFLICT (category) 
	DO UPDATE SET 
		description = EXCLUDED.description, 
		created_at = current_timestamp
	RETURNING id, created_at`

	var dest struct {
		ID        uuid.UUID `db:"id"`
		CreatedAt time.Time `db:"created_at"`
	}

	if err := r.db.GetContext(ctx, &dest, query, category.ID, category.Name, category.Description); err != nil {
		return nil, fmt.Errorf("upsert category in db: %w", err)
	}
	category.ID = dest.ID
	category.CreatedAt = dest.CreatedAt

	return category, nil
}
