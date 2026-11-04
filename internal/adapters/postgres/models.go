package postgres

import (
	"time"

	"github.com/google/uuid"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
)

type userRow struct {
	ID        uuid.UUID `db:"id"`
	Email     string    `db:"email"`
	Role      string    `db:"role"`
	CreatedAt time.Time `db:"created_at"`
}

func (r userRow) toDomain() domain.User {
	return domain.User{
		ID:        r.ID,
		Email:     r.Email,
		Role:      r.Role,
		CreatedAt: r.CreatedAt,
	}
}

type categoryRow struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"category"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
}

func (r categoryRow) toDomain() domain.Category {
	return domain.Category{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		CreatedAt:   r.CreatedAt,
	}
}

type assistantRow struct {
	ID           uuid.UUID `db:"id"`
	Name         string    `db:"name"`
	Description  string    `db:"description"`
	CategoryID   uuid.UUID `db:"category_id"`
	Model        string    `db:"model"`
	SystemPrompt string    `db:"system_prompt"`
	ExUserPrompt string    `db:"ex_user_prompt"`
	IsActivated  bool      `db:"is_activated"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func (r assistantRow) toDomain() domain.Assistant {
	return domain.Assistant{
		ID:           r.ID,
		Name:         r.Name,
		Description:  r.Description,
		CategoryID:   r.CategoryID,
		Model:        r.Model,
		SystemPrompt: r.SystemPrompt,
		ExUserPrompt: r.ExUserPrompt,
		IsActivated:  r.IsActivated,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}

type assistantByIdRow struct {
	assistantRow
	CategoryName string `db:"category"`
}

func (r assistantByIdRow) toDomain() domain.AssistantById {
	return domain.AssistantById{
		Assistant:    r.assistantRow.toDomain(),
		CategoryName: r.CategoryName,
	}
}

type assistantsListRow struct {
	assistantRow
	CategoryName string `db:"category"`
}

func (r assistantsListRow) toDomain() domain.AssistantsList {
	return domain.AssistantsList{
		Assistant:    r.assistantRow.toDomain(),
		CategoryName: r.CategoryName,
	}
}

type assistantRunRow struct {
	ID          uuid.UUID        `db:"id"`
	AssistantID uuid.UUID        `db:"assistant_id"`
	UserID      uuid.UUID        `db:"user_id"`
	UserPrompt  string           `db:"user_prompt"`
	Result      string           `db:"result"`
	Status      domain.RunStatus `db:"status"`
	Error       string           `db:"error"`
	CreatedAt   time.Time        `db:"created_at"`
}

func (r assistantRunRow) toDomain() domain.AssistantRun {
	return domain.AssistantRun{
		ID:          r.ID,
		AssistantID: r.AssistantID,
		UserID:      r.UserID,
		UserPrompt:  r.UserPrompt,
		Result:      r.Result,
		Status:      r.Status,
		Error:       r.Error,
		CreatedAt:   r.CreatedAt,
	}
}

type assistantRunsListRow struct {
	assistantRunRow
	AssistantName string    `db:"assistant_name"`
	CategoryID    uuid.UUID `db:"category_id"`
	CategoryName  string    `db:"category_name"`
	Model         string    `db:"model"`
}

func (r assistantRunsListRow) toDomain() domain.AssistantRunsList {
	return domain.AssistantRunsList{
		AssistantRun:  r.assistantRunRow.toDomain(),
		AssistantName: r.AssistantName,
		CategoryID:    r.CategoryID,
		CategoryName:  r.CategoryName,
		Model:         r.Model,
	}
}
