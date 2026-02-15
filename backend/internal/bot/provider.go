package bot

import (
	"context"
	"errors"
)

var (
	// ErrProviderNotConfigured is returned when bot calls are disabled.
	ErrProviderNotConfigured = errors.New("bot provider not configured")
	// ErrNotImplemented is returned for boilerplate methods not implemented yet.
	ErrNotImplemented = errors.New("bot provider method not implemented")
)

// Provider abstracts bot decision and optional chat generation.
type Provider interface {
	DecideTurn(ctx context.Context, req DecisionRequest) (DecisionResponse, error)
	GenerateChat(ctx context.Context, req ChatRequest) (ChatResponse, error)
}

type noopProvider struct{}

func NewNoopProvider() Provider {
	return &noopProvider{}
}

func (p *noopProvider) DecideTurn(_ context.Context, _ DecisionRequest) (DecisionResponse, error) {
	return DecisionResponse{}, ErrProviderNotConfigured
}

func (p *noopProvider) GenerateChat(_ context.Context, _ ChatRequest) (ChatResponse, error) {
	return ChatResponse{}, ErrProviderNotConfigured
}
