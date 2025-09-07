package models

type GamePhase string

const (
	Setup      	GamePhase = "setup"
	Election   	GamePhase = "election"
	Legislation GamePhase = "legislation"
	Executive  	GamePhase = "executive"
	GameOver   	GamePhase = "game_over"
	Paused     	GamePhase = "paused"
)