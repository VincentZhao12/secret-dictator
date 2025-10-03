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
	Players             []Player       `json:"players"`
	PlayerIndexMap      map[string]int `json:"-"`
	Deck                []Card         `json:"deck"`
	Discard             []Card         `json:"discard"`
	Board               Board          `json:"board"`
	PresidentIndex      int            `json:"president_index"`
	ChancellorIndex     int            `json:"chancellor_index"`
	PrevPresidentIndex  int            `json:"prev_president_index"`
	PrevChancellorIndex int            `json:"prev_chancellor_index"`
	NomineeIndex        int            `json:"nominee_index"`
	Phase               GamePhase      `json:"phase"`
	Votes               []VoteResult   `json:"votes,omitempty"`
	PendingAction       *Action        `json:"pending_action,omitempty"`
	PeekedCards         []Card         `json:"peeked_cards,omitempty"`
	PeekerIndex         int            `json:"peeker_index,omitempty"`
	ResumeOrderIndex    int            `json:"resume_order_index,omitempty"` // Post special election
	ResumePhase         GamePhase      `json:"resume_phase,omitempty"`
	Winner              Team           `json:"winner,omitempty"`
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
		Players:             []Player{},
		PlayerIndexMap:      make(map[string]int),
		Deck:                []Card{},
		Discard:             []Card{},
		Board:               Board{},
		PresidentIndex:      -1,
		ChancellorIndex:     -1,
		PrevPresidentIndex:  -1,
		PrevChancellorIndex: -1,
		NomineeIndex:        -1,
		Phase:               Setup,
		Votes:               nil,
		PendingAction:       nil,
		PeekedCards:         nil,
		PeekerIndex:         -1,
		ResumeOrderIndex:    -1,
		Winner:              TeamUnassigned,
	}
}

func (state *GameState) AddPlayer(id string, username string) (*Player, error) {
	if _, exists := state.PlayerIndexMap[id]; exists {
		return nil, repository.ErrPlayerAlreadyExists
	}

	if len(state.Players) >= 10 {
		return nil, repository.ErrGameFull
	}

	if state.Phase != Setup && (state.Phase == Paused && state.ResumePhase == Setup) {
		return nil, repository.ErrGameInProgress
	}

	state.PlayerIndexMap[id] = len(state.Players)
	player := NewPlayer(id, username)
	state.Players = append(state.Players, player)
	return &player, nil
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
	if forPlayer.ID != obfuscatedState.Players[obfuscatedState.PeekerIndex].ID {
		obfuscatedState.PeekedCards = nil
	}

	// Obfuscate player information
	for i, player := range obfuscatedState.Players {
		if forPlayer.Role != RoleFascist && player.ID != forPlayer.ID {
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
	if obfuscatedState.PeekerIndex != -1 && p.ID != obfuscatedState.Players[obfuscatedState.PeekerIndex].ID {
		obfuscatedState.PeekedCards = nil
	}

	// Obfuscate player information
	for i, player := range obfuscatedState.Players {
		if p.Role != RoleFascist && player.ID != p.ID {
			obfuscatedState.Players[i].Role = RoleHidden
		}

		if player.ID != p.ID {
			player.ID = ""
		}
	}

	return obfuscatedState
}

func (state *GameState) EndGame(winner Team) {
	state.Phase = GameOver
	state.Winner = winner
}

func (state *GameState) EndGameIfNecessary() bool {
	for _, player := range state.Players {
		if player.Role == RoleHitler && player.IsExecuted {
			state.EndGame(TeamLiberal)
			return true
		}
	}

	if state.Board.FascistPolicies >= state.Board.FascistSlots {
		state.EndGame(TeamFascist)
		return true
	}

	if state.Board.LiberalPolicies >= state.Board.LiberalSlots {
		state.EndGame(TeamLiberal)
		return true
	}

	return false
}

func (state *GameState) NewTurn() {
	state.Phase = Nomination
	state.PrevPresidentIndex = state.PresidentIndex
	state.PrevChancellorIndex = state.ChancellorIndex

	if state.ResumeOrderIndex != -1 {
		state.PresidentIndex = state.ResumeOrderIndex
		state.ResumeOrderIndex = -1
	} else {
		nextPresidentIndex := (state.PresidentIndex + 1) % len(state.Players)
		for state.Players[nextPresidentIndex].IsExecuted {
			nextPresidentIndex = (nextPresidentIndex + 1) % len(state.Players)
		}
		state.PresidentIndex = nextPresidentIndex
	}

	state.ChancellorIndex = -1
	state.NomineeIndex = -1
	state.Votes = nil
	state.PendingAction = nil
	state.PeekedCards = nil
	state.PeekerIndex = -1
}

func (state *GameState) PlaceCard(card Card) bool {
	switch card {
	case CardFascist:
		state.Board.FascistPolicies++
	case CardLiberal:
		state.Board.LiberalPolicies++
	}

	if state.Board.ExecutiveActions[state.Board.FascistPolicies] != ActionNone {
		state.Phase = Executive
		action := state.Board.ExecutiveActions[state.Board.FascistPolicies]
		state.PendingAction = &action
	}

	return state.EndGameIfNecessary()
}
