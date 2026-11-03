package postgres

import (
	"context"
	"fmt"

	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
)

func (r *PostgresRepository) RunAssistant(ctx context.Context, runs *domain.AssistantRun) error {
	query := `
	INSERT INTO assistant_run 
		(id, assistant_id, user_id, user_prompt, result, status, error)
	VALUES 
		($1, $2, $3, $4, $5, $6, $7) 
	ON CONFLICT (id) 
	DO UPDATE SET 
		result = EXCLUDED.result, 
		status = EXCLUDED.status, 
		error = EXCLUDED.error`

	if _, err := r.db.ExecContext(ctx, query, runs.ID, runs.AssistantID, runs.UserID, runs.UserPrompt, runs.Result, runs.Status, runs.Error); err != nil {
		return fmt.Errorf("upserting assistant run into db: %w", err)
	}
	return nil
}
