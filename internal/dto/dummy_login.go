package dto

import (
	"github.com/google/uuid"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
)

type DummyLoginInput struct {
	Role string `json:"role"`
}

type DummyLoginOutput struct {
	Token string `json:"token"`
	User  struct {
		ID    uuid.UUID `json:"id"`
		Email string    `json:"email"`
		Role  string    `json:"role"`
	} `json:"user"`
}

func ToDummyLoginDTO(token string, user *domain.User) *DummyLoginOutput {
	return &DummyLoginOutput{
		Token: token,
		User: struct {
			ID    uuid.UUID `json:"id"`
			Email string    `json:"email"`
			Role  string    `json:"role"`
		}{
			ID:    user.ID,
			Email: user.Email,
			Role:  user.Role,
		},
	}
}
