package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
	"github.com/nickchervov/go-ai-agents-catalog/internal/dto"
)

func (s *AgentService) DummyLogin(ctx context.Context, input dto.DummyLoginInput) (*dto.DummyLoginOutput, error) {
	if input.Role != "admin" && input.Role != "user" {
		return nil, domain.ErrIncorrectRole
	}

	val, ok := domain.RoleDatas[input.Role]
	if !ok {
		return nil, domain.ErrIncorrectRole
	}

	id, err := uuid.Parse(val.UserID)
	if err != nil {
		return nil, fmt.Errorf("parsing uuid: %w", err)
	}

	user, err := s.repo.UpsertUser(ctx, &domain.User{ID: id, Email: val.Email, Role: input.Role})
	if err != nil {
		return nil, fmt.Errorf("upserting user: %w", err)
	}

	token, err := s.tokenCreator.CreateJWTToken(user.ID.String(), user.Role, user.Email)
	if err != nil {
		return nil, fmt.Errorf("creating token: %w", err)
	}

	output := dto.ToDummyLoginDTO(token, user)
	return output, nil
}
