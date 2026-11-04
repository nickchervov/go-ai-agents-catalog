package domain

import (
	"time"

	"github.com/google/uuid"
)

type RunStatus string

const (
	Pending RunStatus = "pending"
	Success RunStatus = "success"
	Failed  RunStatus = "failed"
)

type User struct {
	ID        uuid.UUID
	Email     string
	Role      string
	CreatedAt time.Time
}

type Category struct {
	ID          uuid.UUID
	Name        string
	Description string
	CreatedAt   time.Time
}

type Assistant struct {
	ID           uuid.UUID
	Name         string
	Description  string
	CategoryID   uuid.UUID
	Model        string
	SystemPrompt string
	ExUserPrompt string
	IsActivated  bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
type AssistantById struct {
	Assistant
	CategoryName string
}

type AssistantsList struct {
	Assistant
	CategoryName string
}

type AssistantRun struct {
	ID          uuid.UUID
	AssistantID uuid.UUID
	UserID      uuid.UUID
	UserPrompt  string
	Result      string
	Status      RunStatus
	Error       string
	CreatedAt   time.Time
}

type AssistantRunsList struct {
	AssistantRun
	AssistantName string
	CategoryID    uuid.UUID
	CategoryName  string
	Model         string
}
