package bot

import (
	"context"

	"github.com/openai/openai-go/v3"
)

// OpenRouterProvider is a Provider backed by OpenRouter via openai-go.
type OpenRouterProvider struct {
	client       *openai.Client
	defaultModel string
}

func NewOpenRouterProvider(client *openai.Client, defaultModel string) Provider {
	if client == nil {
		return NewNoopProvider()
	}
	return &OpenRouterProvider{
		client:       client,
		defaultModel: defaultModel,
	}
}

func (p *OpenRouterProvider) DecideTurn(ctx context.Context, req DecisionRequest) (DecisionResponse, error) {
	
	return DecisionResponse{}, ErrNotImplemented
}

func (p *OpenRouterProvider) GenerateChat(ctx context.Context, req ChatRequest) (ChatResponse, error) {
	return ChatResponse{}, ErrNotImplemented
}
