package handlers

import (
	"encoding/json"
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
			return
		}
		defer conn.Close()
		queryParams := r.URL.Query()
		gameId := queryParams.Get("game")

		if gameId == "" {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		game, exists := Manager.Games[gameId]
		if !exists || game == nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		playerId := queryParams.Get("player")
		err = game.AddConnection(playerId, conn)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
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
				return
			}

			var message messages.Message
			err = json.Unmarshal(messageBytes, &message)
			if err != nil {
				println("Malformed web socket message")
			}

			switch message.GetType() {
			case messages.MesasgeTypeAction:
				message := message.(*messages.ActionMessage)
				game.ActionChan <- *message
			}
		}
	}
}
