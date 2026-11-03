package llm

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type LLMAgent struct {
	client *openai.Client
}

func NewLLMAgent(api, baseURL string) *LLMAgent {
	config := openai.DefaultConfig(api)
	config.BaseURL = baseURL

	client := openai.NewClientWithConfig(config)

	return &LLMAgent{client: client}
}

func (a *LLMAgent) SendMessage(ctx context.Context, model, systemPrompt, userPrompt string) (string, error) {
	resp, err := a.client.CreateResponse(ctx, openai.CreateResponseRequest{
		Model:        model,
		Instructions: systemPrompt,
		Input:        userPrompt,
	})
	if err != nil {
		return "", fmt.Errorf("creating response from llm: %w", err)
	}

	return resp.GetOutputText(), nil
}
