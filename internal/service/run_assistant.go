package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
	"github.com/nickchervov/go-ai-agents-catalog/internal/dto"
)

func (s *AgentService) RunAssistant(ctx context.Context, input dto.RunAssistantInput) (*dto.RunAssistantOutput, error) {
	if input.UserPrompt == "" {
		return nil, domain.ErrEmptyUserPrompt
	}

	id := uuid.New()

	runs := domain.AssistantRun{
		ID:          id,
		AssistantID: input.AssistantId,
		UserID:      input.UserID,
		UserPrompt:  input.UserPrompt,
		Status:      domain.Pending,
	}

	assistant, err := s.repo.AssistantByID(ctx, input.AssistantId)
	if err != nil {
		return nil, fmt.Errorf("getting assistant by id: %w", err)
	}

	if !assistant.IsActivated {
		return nil, domain.ErrAssistantInactive
	}

	if err := s.repo.RunAssistant(ctx, &runs); err != nil {
		return nil, fmt.Errorf("upsert run assistant: %w", err)
	}

	result, err := s.llm.SendMessage(ctx, assistant.Model, assistant.SystemPrompt, runs.UserPrompt)
	if err != nil {
		runs.Error = err.Error()
		runs.Status = domain.Failed

		if err := s.repo.RunAssistant(ctx, &runs); err != nil {
			return nil, fmt.Errorf("upsert run assistant after err: %w", err)
		}

		return nil, domain.ErrLLMProvider
	}

	runs.Result = result
	runs.Status = domain.Success

	if err := s.repo.RunAssistant(ctx, &runs); err != nil {
		return nil, fmt.Errorf("upsert run assistant after success: %w", err)
	}

	output := dto.RunAssistantOutput{
		ID:            id,
		AssistantID:   runs.AssistantID,
		AssistantName: assistant.Name,
		CategoryID:    assistant.CategoryID,
		UserID:        runs.UserID,
		Model:         assistant.Model,
		UserPrompt:    runs.UserPrompt,
		Output:        runs.Result,
		Status:        runs.Status,
		Error:         runs.Error,
	}
	return &output, nil
}
