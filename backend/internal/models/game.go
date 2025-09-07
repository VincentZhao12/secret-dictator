package models

type Game struct {
	ID      			string `json:"id"`
	Host    			string `json:"host"`
	Players 			[]Player `json:"players"`
	Deck    			[]Card `json:"deck"`
	Discard 			[]Card `json:"discard"`
	Board   			Board `json:"board"`
	CurrentPresident 	string `json:"current_president"`
	CurrentChancellor 	string `json:"current_chancellor"`
	PreviousPresident 	string `json:"previous_president"`
	PreviousChancellor 	string `json:"previous_chancellor"`
	PresidentIndex 		int `json:"president_index"`
	Nominee 			string `json:"nominee"`
	Phase            	GamePhase `json:"phase"`
	Votes            	map[string]int `json:"votes"`
	PendingAction      	*Action         `json:"pending_action,omitempty"`
}