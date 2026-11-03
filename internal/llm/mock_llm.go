package llm

import (
	"context"
	"fmt"
)

type MockLLM struct {
	WantErr bool
}

func (m *MockLLM) SendMessage(ctx context.Context, model, systemPrompt, userPrompt string) (string, error) {
	if m.WantErr {
		return "", fmt.Errorf("creating response from mock llm")
	}
	return "Nice prompt take result back", nil
}
