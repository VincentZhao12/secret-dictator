package bot

import "context"

// Orchestrator coordinates game-triggered bot calls with the Provider.
type Orchestrator struct {
	provider Provider
}

func NewOrchestrator(provider Provider) *Orchestrator {
	if provider == nil {
		provider = NewNoopProvider()
	}
	return &Orchestrator{
		provider: provider,
	}
}

func (o *Orchestrator) Provider() Provider {
	return o.provider
}

func (o *Orchestrator) DecideTurn(ctx context.Context, req DecisionRequest) (DecisionResponse, error) {
	return o.provider.DecideTurn(ctx, req)
}

func (o *Orchestrator) GenerateChat(ctx context.Context, req ChatRequest) (ChatResponse, error) {
	return o.provider.GenerateChat(ctx, req)
}
