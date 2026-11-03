package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
)

func (r *PostgresRepository) UpdateAssistantByID(ctx context.Context, id uuid.UUID, assistant *domain.Assistant) (*domain.AssistantsResponse, error) {
	if assistant.IsActivated == nil {
		t := true
		assistant.IsActivated = &t
	}
	var newAssistant domain.AssistantsResponse
	query := `
	WITH updated_assistant AS (
		UPDATE assistants 
		SET 
			name = $1, 
			description = $2, 
			category_id = $3, 
			model = $4, 
			system_prompt = $5, 
			ex_user_prompt = $6, 
			is_activated = $7,
			updated_at = $8
		WHERE 
			id = $9
		RETURNING
			id,
			name,
			description,
			category_id,
			model,
			system_prompt, 
			ex_user_prompt, 
			is_activated, 
			updated_at,
			created_at
	)
			
	SELECT
		ua.id,
		ua.name,
		ua.description,
		ua.category_id,
		c.category as category,
		ua.model,
		ua.system_prompt, 
		ua.ex_user_prompt, 
		ua.is_activated, 
		ua.updated_at,
		ua.created_at
	FROM updated_assistant as ua
	LEFT JOIN categories AS c ON ua.category_id = c.id`

	err := r.db.GetContext(ctx, &newAssistant, query, assistant.Name, assistant.Description,
		assistant.CategoryID, assistant.Model, assistant.SystemPrompt, assistant.ExUserPrompt, assistant.IsActivated, assistant.UpdatedAt, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrAssistantNotFound
		}
		return nil, fmt.Errorf("inserting new assistant by id in db: %w", err)
	}

	return &newAssistant, nil
}
