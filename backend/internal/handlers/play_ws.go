package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/VincentZhao12/secret-hitler/backend/internal/game"
	"github.com/VincentZhao12/secret-hitler/backend/internal/messages"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // In production, verify origin
}

func scheduleDeleteGame(m *game.Manager, game *game.Game) {
	go func() {
		time.Sleep(10 * time.Minute)
		if g, exists := m.GetGame(game.ID); exists && g.CanBeDeleted() {
			m.RemoveGame(game.ID)
			fmt.Println("Deleted game", game.ID, "due to inactivity")
		}
	}()
}

func Play(Manager *game.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			// TODO: Log failed connections
			conn.WriteJSON(messages.NewConnectionErrorMessage("server", "Failed to upgrade to WebSocket: "+err.Error(), messages.ConnectionErrorTypeGameInvalid))
			return
		}
		defer conn.Close()
		queryParams := r.URL.Query()
		gameId := queryParams.Get("game")

		if gameId == "" {
			conn.WriteJSON(messages.NewConnectionErrorMessage("server", "Missing game ID in query parameters", messages.ConnectionErrorTypeGameInvalid))
			fmt.Println("no game id")
			return
		}

		game, exists := Manager.GetGame(gameId)
		if !exists || game == nil {
			conn.WriteJSON(messages.NewConnectionErrorMessage("server", "Game not found", messages.ConnectionErrorTypeGameInvalid))
			fmt.Println("no game found")
			return
		}

		playerId := queryParams.Get("player")
		err = game.AddConnection(playerId, conn)
		if err != nil {
			conn.WriteJSON(messages.NewConnectionErrorMessage("server", "Failed to add connection: "+err.Error(), messages.ConnectionErrorTypePlayerInvalid))
			fmt.Println("no player found")
			return
		}

		defer func() {
			game.DropConnection(playerId)
			scheduleDeleteGame(Manager, game)
		}()

		for {
			_, messageBytes, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("error reading message")
				return
			}

			var envelope struct {
				Base messages.BaseMessage `json:"base_message"`
			}
			if err := json.Unmarshal(messageBytes, &envelope); err != nil {
				fmt.Println("Malformed web socket message: could not parse envelope")
				continue
			}

			switch envelope.Base.GetType() {
			case messages.MessageTypeAction:
				var action messages.ActionMessage
				if err := json.Unmarshal(messageBytes, &action); err != nil {
					fmt.Println("Malformed action message")
					continue
				}
				// Never trust client-provided sender IDs; bind actions to this socket's player.
				action.SenderID = playerId
				game.ActionChan <- action
			default:
				fmt.Println("Unexpected message type")
			}
		}
	}
}
