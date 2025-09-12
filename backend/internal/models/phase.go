package models

type GamePhase string

const (
	Setup      		GamePhase = "setup"
	Nomination 		GamePhase = "nomination"
	Election   		GamePhase = "election"
	Legislation 	GamePhase = "legislation"
	Executive  		GamePhase = "executive"
	GameOver   		GamePhase = "game_over"
	Paused     		GamePhase = "paused"
)