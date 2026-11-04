package mocks

import (
	"context"
	"fmt"
)

type MockLLM struct {
	Result string
	Err    error
}

func (m *MockLLM) SendMessage(ctx context.Context, model, systemPrompt, userPrompt string) (string, error) {
	if m.Err != nil {
		return "", fmt.Errorf("creating response from mock llm: %w", m.Err)
	}
	return "Nice prompt take result back", nil
}
