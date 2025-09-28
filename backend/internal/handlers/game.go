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

type JoinGameRequest struct {
	GameID string `json:"game_id"`
	Username string `json:"username"` 
}

type JoinGameResponse struct {
	GameID string `json:"game_id"`
	PlayerID string `json:"player_id"` 
}

func JoinGame(Manager *game.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req JoinGameRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		game, exists := Manager.GetGame(req.GameID)
		if game == nil  || !exists {
			http.Error(w, "Invalid game id", http.StatusBadRequest)
		}
		
		player, err := game.NewPlayer(req.Username);
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp := JoinGameResponse{
			GameID: game.ID,
			PlayerID: player.ID,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated) // sets 201
		json.NewEncoder(w).Encode(resp)   // encodes and writes JSON
	}
}