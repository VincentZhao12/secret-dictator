package messages

import "github.com/VincentZhao12/secret-hitler/backend/internal/models"

const (
	MessageTypeGameState MessageType = "game_state"
)

type GameStateMessage struct {
	BaseMessage
	GameState models.GameState `json:"game_state" tstype:"GameState"`
}

func NewGameStateMessage(senderID string, gameState models.GameState) *GameStateMessage {
	return &GameStateMessage{
		BaseMessage: BaseMessage{
			Type:     MessageTypeGameState,
			SenderID: senderID,
		},
		GameState: gameState,
	}
}