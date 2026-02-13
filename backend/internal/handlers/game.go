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
		newGame := game.NewGame(Manager)
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
	GameID   string `json:"game_id"`
	Username string `json:"username"`
}

type JoinGameResponse struct {
	GameID   string `json:"game_id"`
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
		if game == nil || !exists {
			http.Error(w, "Invalid game id", http.StatusBadRequest)
			return
		}

		player, err := game.NewPlayer(req.Username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp := JoinGameResponse{
			GameID:   game.ID,
			PlayerID: player.ID,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated) // sets 201
		json.NewEncoder(w).Encode(resp)   // encodes and writes JSON
	}
}

type AddBotRequest struct {
	GameID    string `json:"game_id"`
	HostID    string `json:"host_id"`
	ModelSlug string `json:"model_slug"`
}

type AddBotResponse struct {
	GameID   string `json:"game_id"`
	PlayerID string `json:"player_id"`
}

func AddBot(Manager *game.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AddBotRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		game, exists := Manager.GetGame(req.GameID)
		if game == nil || !exists {
			http.Error(w, "Invalid game id", http.StatusBadRequest)
			return
		}

		if req.HostID != game.HostID {
			http.Error(w, "Only host can manage bots", http.StatusForbidden)
			return
		}

		modelSlug := req.ModelSlug
		if modelSlug == "" {
			modelSlug = "openai/gpt-oss-120b:free"
		}

		player, err := game.AddBotForHost(req.HostID, modelSlug)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp := AddBotResponse{
			GameID:   game.ID,
			PlayerID: player.ID,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	}
}

type RemoveBotRequest struct {
	GameID      string `json:"game_id"`
	HostID      string `json:"host_id"`
	BotPlayerID string `json:"bot_player_id"`
}

func RemoveBot(Manager *game.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RemoveBotRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		game, exists := Manager.GetGame(req.GameID)
		if game == nil || !exists {
			http.Error(w, "Invalid game id", http.StatusBadRequest)
			return
		}

		if req.HostID != game.HostID {
			http.Error(w, "Only host can manage bots", http.StatusForbidden)
			return
		}

		if err := game.RemoveBotForHost(req.HostID, req.BotPlayerID); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
