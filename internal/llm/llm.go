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

func (a *LLMAgent) Models(ctx context.Context) ([]string, error) {
	models, err := a.client.ListModels(ctx)
	if err != nil {
		return nil, fmt.Errorf("get list models: %w", err)
	}

	modelsIDs := make([]string, 0, len(models.Models))

	for _, m := range models.Models {
		modelsIDs = append(modelsIDs, m.ID)
	}

	return modelsIDs, nil
}
