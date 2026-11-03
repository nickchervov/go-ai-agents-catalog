package postgres

import (
	"context"
	"fmt"

	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
)

func (r *PostgresRepository) UpsertUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	var upsertedUser domain.User
	query := `
	INSERT INTO users 
		(id, email, role) 
	VALUES 
		($1, $2, $3)
	ON CONFLICT (id) 
	DO UPDATE SET 
		email = EXCLUDED.email
	RETURNING id, email, role, created_at`

	if err := r.db.GetContext(ctx, &upsertedUser, query, user.ID, user.Email, user.Role); err != nil {
		return nil, fmt.Errorf("upserting user in db: %w", err)
	}

	return &upsertedUser, nil
}
