package models

type GamePhase string

const (
	Setup      		GamePhase = "setup"
	Nomination 		GamePhase = "nomination"
	Election   		GamePhase = "election"
	Legislation1 	GamePhase = "legislation"
	Legislation2 	GamePhase = "legislation"
	Executive  		GamePhase = "executive"
	GameOver   		GamePhase = "game_over"
	Paused     		GamePhase = "paused"
)