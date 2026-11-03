package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `db:"id"`
	Email     string    `db:"email"`
	Role      string    `db:"role"`
	CreatedAt time.Time `db:"created_at"`
}

type Category struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"category"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
}

type Assistant struct {
	ID           uuid.UUID `db:"id"`
	Name         string    `db:"name"`
	Description  string    `db:"description"`
	CategoryID   uuid.UUID `db:"category_id"`
	Model        string    `db:"model"`
	SystemPrompt string    `db:"system_prompt"`
	ExUserPrompt string    `db:"ex_user_prompt"`
	IsActivated  *bool     `db:"is_activated"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
type AssistantByIdResponse struct {
	ID           uuid.UUID `db:"id"`
	Name         string    `db:"name"`
	Description  string    `db:"description"`
	CategoryID   uuid.UUID `db:"category_id"`
	CategoryName string    `db:"category"`
	Model        string    `db:"model"`
	SystemPrompt string    `db:"system_prompt"`
	ExUserPrompt string    `db:"ex_user_prompt"`
	IsActivated  bool      `db:"is_activated"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type AssistantsResponse struct {
	ID           uuid.UUID `db:"id"`
	Name         string    `db:"name"`
	Description  string    `db:"description"`
	CategoryID   uuid.UUID `db:"category_id"`
	CategoryName string    `db:"category"`
	Model        string    `db:"model"`
	SystemPrompt string    `db:"system_prompt"`
	ExUserPrompt string    `db:"ex_user_prompt"`
	IsActivated  bool      `db:"is_activated"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
	TotalCount   int       `db:"total_count"`
}

type RunStatus string

const (
	Pending RunStatus = "pending"
	Success RunStatus = "success"
	Failed  RunStatus = "failed"
)

type AssistantRun struct {
	ID          uuid.UUID `db:"id"`
	AssistantID uuid.UUID `db:"assistant_id"`
	UserID      uuid.UUID `db:"user_id"`
	UserPrompt  string    `db:"user_prompt"`
	Result      string    `db:"result"`
	Status      RunStatus `db:"status"`
	Error       string    `db:"error"`
	CreatedAt   time.Time `db:"created_at"`
}

type UserRunsResponse struct {
	ID            uuid.UUID `db:"id"`
	AssistantID   uuid.UUID `db:"assistant_id"`
	AssistantName string    `db:"assistant_name"`
	CategoryID    uuid.UUID `db:"category_id"`
	CategoryName  string    `db:"category_name"`
	UserID        uuid.UUID `db:"user_id"`
	Model         string    `db:"model"`
	UserPrompt    string    `db:"user_prompt"`
	Output        string    `db:"output"`
	Status        string    `db:"status"`
	Error         string    `db:"error"`
	CreatedAt     time.Time `db:"created_at"`
	Total         int       `db:"total"`
}

type AdminRunsResponse struct {
	ID            uuid.UUID `db:"id"`
	AssistantID   uuid.UUID `db:"assistant_id"`
	AssistantName string    `db:"assistant_name"`
	CategoryID    uuid.UUID `db:"category_id"`
	CategoryName  string    `db:"category_name"`
	UserID        uuid.UUID `db:"user_id"`
	Model         string    `db:"model"`
	UserPrompt    string    `db:"user_prompt"`
	Output        string    `db:"output"`
	Status        string    `db:"status"`
	Error         string    `db:"error"`
	CreatedAt     time.Time `db:"created_at"`
	Total         int       `db:"total"`
}
