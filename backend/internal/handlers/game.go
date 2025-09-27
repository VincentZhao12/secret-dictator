package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/VincentZhao12/secret-hitler/backend/internal/game"
)

type CreateGameResponse struct {
	GameID string `json:"game_id"`
}


func CreateGame(Manager *game.Manager) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		newGame := game.NewGame()
		gameID := Manager.AddGame(newGame)

		resp := CreateGameResponse{
			GameID: gameID,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated) // sets 201
		json.NewEncoder(w).Encode(resp)   // encodes and writes JSON
	}
}
