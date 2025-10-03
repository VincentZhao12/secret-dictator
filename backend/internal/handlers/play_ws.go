package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/VincentZhao12/secret-hitler/backend/internal/game"
	"github.com/VincentZhao12/secret-hitler/backend/internal/messages"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // In production, verify origin
}

func Play(Manager *game.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			// TODO: Log failed connections
			fmt.Println("failed to upgrade to web socket connection")
			return
		}
		defer conn.Close()
		queryParams := r.URL.Query()
		gameId := queryParams.Get("game")

		if gameId == "" {
			// http.Error(w, "Invalid request payload", http.StatusBadRequest)
			fmt.Println("no game id")
			return
		}

		game, exists := Manager.GetGame(gameId)
		if !exists || game == nil {
			// http.Error(w, "Invalid request payload", http.StatusBadRequest)
			fmt.Println("no game found")
			return
		}

		playerId := queryParams.Get("player")
		err = game.AddConnection(playerId, conn)
		if err != nil {
			fmt.Println("no player found")
			// http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		defer func() {
			game.DropConnection(playerId)
			if game.CanBeDeleted() {
				Manager.RemoveGame(gameId)
			}
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
				game.ActionChan <- action
			default:
				fmt.Println("Unexpected message type")
			}
		}
	}
}
