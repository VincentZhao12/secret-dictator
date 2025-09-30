package models

type GamePhase string

const (
	Setup        GamePhase = "setup" // TODO
	Nomination   GamePhase = "nomination"
	Election     GamePhase = "election"
	Legislation1 GamePhase = "legislation1"
	Legislation2 GamePhase = "legislation2"
	Executive    GamePhase = "executive"
	GameOver     GamePhase = "game_over"
	Paused       GamePhase = "paused"
)
