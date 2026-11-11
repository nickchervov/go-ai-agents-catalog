package mocks

import (
	"context"
	"fmt"
)

type MockLLM struct {
	Counter int
	Result  string
	Err     error
}

func (m *MockLLM) SendMessage(ctx context.Context, model, systemPrompt, userPrompt string) (string, error) {
	if m.Err != nil {
		return "", fmt.Errorf("creating response from mock llm: %w", m.Err)
	}
	m.Counter++
	return m.Result, nil
}

func (m *MockLLM) Models(ctx context.Context) ([]string, error) {
	if m.Err != nil {
		return nil, fmt.Errorf("getting models list: %w", m.Err)
	}

	return []string{"mock-model"}, nil
}
