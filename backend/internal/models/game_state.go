package models

import (
	"math/rand"

	"github.com/VincentZhao12/secret-hitler/backend/internal/repository"
)

type VoteResult int

const (
	VotePending VoteResult = iota
	VoteJa
	VoteNein
)

type GameState struct {
	Players 			[]Player 		`json:"players"`
	PlayerIndexMap		map[string]int 	`json:"-"`
	Deck    			[]Card 			`json:"deck"`
	Discard 			[]Card 			`json:"discard"`
	Board   			Board 			`json:"board"`
	PresidentIndex 		int 			`json:"president_index"`
	ChancellorIndex 	int 			`json:"chancellor_index"`
	PrevPresidentIndex  int 			`json:"prev_president_index"`
	PrevChancellorIndex int 			`json:"prev_chancellor_index"`
	NomineeIndex    	int 			`json:"nominee_index"`
	Phase            	GamePhase 		`json:"phase"`
	Votes            	[]VoteResult 	`json:"votes,omitempty"`
	PendingAction      	*Action     	`json:"pending_action,omitempty"`
	PeekedCards        	[]Card      	`json:"peeked_cards,omitempty"`
	PeekerIndex			int      		`json:"peeker_index,omitempty"`
	ExecutedPlayers    	[]string    	`json:"executed_players"`
	ResumeOrderIndex	int 			`json:"resume_order_index,omitempty"`
}

func createDeck() []Card {
	deck := make([]Card, 17)
	for i := 0; i < 6; i++ {
		deck[i] = CardLiberal
	}
	for i := 6; i < 17; i++ {
		deck[i] = CardFascist
	}
	return deck
}

func NewGameState() GameState {
	return GameState{
		Players:       		[]Player{},
		PlayerIndexMap:     make(map[string]int),
		Deck:          		[]Card{},
		Discard:       		[]Card{},
		Board:         		Board{},
		PresidentIndex:     -1,
		ChancellorIndex:    -1,
		PrevPresidentIndex: -1,
		PrevChancellorIndex:-1,
		NomineeIndex:       -1,
		Phase:         		Setup,
		Votes:         		nil,
		PendingAction:      nil,
		PeekedCards:      	nil,
		PeekerIndex:         	-1,
		ExecutedPlayers:    []string{},
	}
}

func (state *GameState) AddPlayer(id string, username string) error {
	if _, exists := state.PlayerIndexMap[id]; exists {
		return repository.ErrPlayerAlreadyExists
	}

	if len(state.Players) >= 10 {
		return repository.ErrGameFull
	}

	if state.Phase != Setup {
		return repository.ErrGameInProgress
	}

	state.PlayerIndexMap[id] = len(state.Players)
	state.Players = append(state.Players, NewPlayer(id, username, RoleUnassigned))
	return nil
}

func (state *GameState) RemovePlayer(id string) error {
	if state.Phase != Setup {
		return repository.ErrGameInProgress
	}

	index, exists := state.PlayerIndexMap[id]
	if !exists {
		return nil
	}

	lastPlayer := state.Players[len(state.Players)-1]
	state.Players[index] = lastPlayer
	state.PlayerIndexMap[lastPlayer.ID] = index
	state.Players = state.Players[:len(state.Players)-1]
	delete(state.PlayerIndexMap, id)

	return nil
}

func (state *GameState) AssignRoles() error {
	numPlayers := len(state.Players)
	var roles []PlayerRole
	
	switch numPlayers {
	case 5:
		roles = []PlayerRole{RoleHitler, RoleFascist, RoleLiberal, RoleLiberal, RoleLiberal}
	case 6:
		roles = []PlayerRole{RoleHitler, RoleFascist, RoleLiberal, RoleLiberal, RoleLiberal, RoleLiberal}
	case 7:
		roles = []PlayerRole{RoleHitler, RoleFascist, RoleFascist, RoleLiberal, RoleLiberal, RoleLiberal, RoleLiberal}
	case 8:
		roles = []PlayerRole{RoleHitler, RoleFascist, RoleFascist, RoleLiberal, RoleLiberal, RoleLiberal, RoleLiberal, RoleLiberal}
	case 9:
		roles = []PlayerRole{RoleHitler, RoleFascist, RoleFascist, RoleLiberal, RoleLiberal, RoleLiberal, RoleLiberal, RoleLiberal, RoleLiberal}
	case 10:
		roles = []PlayerRole{RoleHitler, RoleFascist, RoleFascist, RoleLiberal, RoleLiberal, RoleLiberal, RoleLiberal, RoleLiberal, RoleLiberal, RoleLiberal}
	default:
		return repository.ErrInvalidPlayerCount
	}

	rand.Shuffle(len(roles), func(i, j int) {
		roles[i], roles[j] = roles[j], roles[i]
	})

	for i := range state.Players {
		state.Players[i].Role = roles[i]
	}
	return nil
}

func (state *GameState) StartGame() error {
	if state.Phase != Setup {
		return repository.ErrGameInProgress
	}

	if len(state.Players) < 5 || len(state.Players) > 10 {
		return repository.ErrInvalidPlayerCount
	}
	
	board, err := NewBoard(len(state.Players))
	if err != nil {
		return err
	}
	state.Board = board
	state.Phase = Election
	state.PresidentIndex = rand.Intn(len(state.Players))
	state.Discard = createDeck()
	state.Deck = []Card{}

	error := state.AssignRoles()
	if error != nil {
		return error
	}

	state.ShuffleDeck()

	return nil
}

func (state *GameState) ShuffleDeck() {
	rand.Shuffle(len(state.Discard), func(i, j int) {
		state.Discard[i], state.Discard[j] = state.Discard[j], state.Discard[i]
	})

	state.Deck = append(state.Discard, state.Deck...)
	state.Discard = []Card{}
}

func (state *GameState) RevealPlayer(forPlayer Player, player Player) GameState {
	obfuscatedState := *state

	// Obfuscate deck
	for i := range obfuscatedState.Deck {
		obfuscatedState.Deck[i] = CardHidden
	}

	// Obfuscate discard pile
	for i := range obfuscatedState.Discard {
		obfuscatedState.Discard[i] = CardHidden
	}

	// Obfuscate hands
	if forPlayer.ID != obfuscatedState.Players[obfuscatedState.PeekerIndex].ID{
		obfuscatedState.PeekedCards = nil
	}

	// Obfuscate player information
	for i, player := range obfuscatedState.Players {
		if forPlayer.Role != RoleFascist &&  player.ID != forPlayer.ID {
			obfuscatedState.Players[i].Role = RoleHidden
		}

		player.ID = ""
	}

	return obfuscatedState
}


func (state *GameState) ObfuscateGameState(p Player) GameState {
	obfuscatedState := *state

	// Obfuscate deck
	for i := range obfuscatedState.Deck {
		obfuscatedState.Deck[i] = CardHidden
	}

	// Obfuscate discard pile
	for i := range obfuscatedState.Discard {
		obfuscatedState.Discard[i] = CardHidden
	}

	// Obfuscate hands
	if p.ID != obfuscatedState.Players[obfuscatedState.PeekerIndex].ID{
		obfuscatedState.PeekedCards = nil
	}

	// Obfuscate player information
	for i, player := range obfuscatedState.Players {
		if p.Role != RoleFascist &&  player.ID != p.ID {
			obfuscatedState.Players[i].Role = RoleHidden
		}

		player.ID = ""
	}

	return obfuscatedState
}