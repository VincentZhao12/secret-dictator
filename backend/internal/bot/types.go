package bot

import "github.com/VincentZhao12/secret-hitler/backend/internal/models"

// DecisionRequest is the normalized input passed to a bot provider.
type DecisionRequest struct {
	GameID      string
	BotPlayerID string
	ModelSlug   string
	State       models.GameState
}

// DecisionResponse is a provider-parsed bot action candidate.
type DecisionResponse struct {
	Action      models.Action
	TargetIndex int
	Card        models.Card
	Reasoning   string
}

// ChatRequest is the normalized input for optional bot chat generation.
type ChatRequest struct {
	GameID      string
	BotPlayerID string
	ModelSlug   string
	State       models.GameState
}

// ChatResponse contains a single generated chat line for the bot.
type ChatResponse struct {
	Text string
}
