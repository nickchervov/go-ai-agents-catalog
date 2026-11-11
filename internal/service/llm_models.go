package service

import (
	"context"

	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
	"github.com/nickchervov/go-ai-agents-catalog/internal/dto"
)

func (s *AgentService) Models(ctx context.Context) (*dto.ModelsOutput, error) {
	models, err := s.llm.Models(ctx)
	if err != nil {
		return nil, domain.ErrLLMProvider
	}

	return &dto.ModelsOutput{Models: models}, nil
}
