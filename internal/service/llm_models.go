package service

import (
	"context"
	"time"

	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
	"github.com/nickchervov/go-ai-agents-catalog/internal/dto"
)

func (s *AgentService) Models(ctx context.Context) (*dto.ModelsOutput, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	models, err := s.llm.Models(timeoutCtx)
	if err != nil {
		return nil, domain.ErrLLMProvider
	}

	return &dto.ModelsOutput{Models: models}, nil
}
