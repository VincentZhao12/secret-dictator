package models

import (
	"math/rand"

	"github.com/VincentZhao12/secret-hitler/backend/internal/repository"
)

type VoteResult int

const (
	VotePending VoteResult = iota
	VoteHidden
	VoteJa
	VoteNein
)

type ChatEntry struct {
	SenderID   string `json:"sender_id"`
	SenderName string `json:"sender_name"`
	Text       string `json:"text"`
	SentAtUnix int64  `json:"sent_at_unix"`
}

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
	HostID              string         `json:"host_id"`
	ChatHistory         []ChatEntry    `json:"chat_history"`
	BotNotes            map[string][]string `json:"bot_notes,omitempty"`
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
		HostID:              "",
		ChatHistory:         []ChatEntry{},
		BotNotes:            make(map[string][]string),
	}
}

// GetPlayer safely gets a player at the given index, returning nil if index is out of bounds
func (state *GameState) GetPlayer(index int) *Player {
	if index < 0 || index >= len(state.Players) {
		return nil
	}
	return &state.Players[index]
}

// GetPlayerByID gets a player by their ID, returning nil if not found
func (state *GameState) GetPlayerByID(id string) *Player {
	index, exists := state.PlayerIndexMap[id]
	if !exists {
		return nil
	}
	return state.GetPlayer(index)
}

func (state *GameState) AddPlayer(id string, username string) (*Player, error) {
	if _, exists := state.PlayerIndexMap[id]; exists {
		return nil, repository.ErrPlayerAlreadyExists
	}

	if len(state.Players) >= 10 {
		return nil, repository.ErrGameFull
	}

	if state.Phase != Setup && (state.Phase != Paused || state.ResumePhase != Setup) {
		return nil, repository.ErrGameInProgress
	}

	state.PlayerIndexMap[id] = len(state.Players)
	player := NewPlayer(id, username)
	state.Players = append(state.Players, player)
	return &player, nil
}

func (state *GameState) AddBotPlayer(id string, modelSlug string) (*Player, error) {
	if _, exists := state.PlayerIndexMap[id]; exists {
		return nil, repository.ErrPlayerAlreadyExists
	}

	if len(state.Players) >= 10 {
		return nil, repository.ErrGameFull
	}

	if state.Phase != Setup {
		return nil, repository.ErrGameInProgress
	}

	state.PlayerIndexMap[id] = len(state.Players)
	username := state.nextBotUsername()
	player := NewBotPlayer(id, modelSlug, username)
	state.Players = append(state.Players, player)
	return &player, nil
}

func (state *GameState) nextBotUsername() string {
	var username string = BotUsernames[rand.Intn(len(BotUsernames))]
	for _, player := range state.Players {
		if player.Username == username {
			return state.nextBotUsername()
		}
	}
	return username
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
	state.Phase = Nomination
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

func (state *GameState) RevealPlayer(viewer Player, revealed Player) GameState {
	obfuscatedState := *state

	// Deep copy deck
	obfuscatedState.Deck = append([]Card(nil), state.Deck...)

	// Deep copy discard
	obfuscatedState.Discard = append([]Card(nil), state.Discard...)

	// Deep copy chat history without leaking player IDs
	obfuscatedState.ChatHistory = make([]ChatEntry, len(state.ChatHistory))
	for i, chat := range state.ChatHistory {
		if chat.SenderID != viewer.ID {
			chat.SenderID = ""
		}
		obfuscatedState.ChatHistory[i] = chat
	}

	// Hide peeked cards unless viewer is the peeker
	peeker := obfuscatedState.GetPlayer(obfuscatedState.PeekerIndex)
	if peeker != nil && viewer.ID != peeker.ID {
		obfuscatedState.PeekedCards = nil
	}

	// Obfuscate players
	obfuscatedState.Players = make([]Player, len(state.Players))
	for i := range obfuscatedState.Players {
		p := state.Players[i]

		// Hide roles based on viewer/revealed relationship
		if viewer.Role == RoleLiberal && p.ID != viewer.ID && p.ID != revealed.ID {
			p.Role = RoleHidden
		} else if viewer.Role == RoleFascist && p.Role == RoleLiberal && p.ID != viewer.ID && p.ID != revealed.ID {
			p.Role = RoleHidden
		}

		// Hide IDs of other players (optional)
		if p.ID != viewer.ID {
			p.ID = ""
		}

		obfuscatedState.Players[i] = p
	}

	return obfuscatedState
}

func (state *GameState) ObfuscateGameState(p Player) GameState {
	obfuscatedState := *state

	if p.ID != state.HostID {
		obfuscatedState.HostID = ""
	}

	// Obfuscate deck
	obfuscatedState.Deck = make([]Card, len(state.Deck))
	for i := range obfuscatedState.Deck {
		obfuscatedState.Deck[i] = state.Deck[i]
	}

	// Obfuscate discard pile
	obfuscatedState.Discard = make([]Card, len(state.Discard))
	for i := range obfuscatedState.Discard {
		obfuscatedState.Discard[i] = state.Discard[i]
	}

	// Deep copy chat history without leaking player IDs
	obfuscatedState.ChatHistory = make([]ChatEntry, len(state.ChatHistory))
	for i, chat := range state.ChatHistory {
		if chat.SenderID != p.ID {
			chat.SenderID = ""
		}
		obfuscatedState.ChatHistory[i] = chat
	}

	// Obfuscate hands
	peeker := obfuscatedState.GetPlayer(obfuscatedState.PeekerIndex)
	if obfuscatedState.PeekerIndex == -1 || peeker == nil || p.ID != peeker.ID {
		obfuscatedState.PeekedCards = nil
	}

	// Obfuscate player information
	obfuscatedState.Players = make([]Player, len(state.Players))
	for i := range obfuscatedState.Players {
		player := state.Players[i]
		if p.Role != RoleFascist && p.ID != player.ID {
			player.Role = RoleHidden
		}

		if p.ID != player.ID {
			player.ID = ""
		}
		obfuscatedState.Players[i] = player
	}

	if state.Votes != nil && state.Phase == Election {
		obfuscatedState.Votes = make([]VoteResult, len(state.Votes))
		for i := range obfuscatedState.Votes {
			if state.Players[i].ID == p.ID {
				obfuscatedState.Votes[i] = state.Votes[i]
			} else {
				obfuscatedState.Votes[i] = VoteHidden
			}
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
		nextPlayer := state.GetPlayer(nextPresidentIndex)
		for nextPlayer != nil && nextPlayer.IsExecuted {
			nextPresidentIndex = (nextPresidentIndex + 1) % len(state.Players)
			nextPlayer = state.GetPlayer(nextPresidentIndex)
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

	if action, exists := state.Board.ExecutiveActions[state.Board.FascistPolicies]; card == CardFascist && exists && action != ActionNone {
		state.Phase = Executive
		action := state.Board.ExecutiveActions[state.Board.FascistPolicies]
		state.PendingAction = &action

		if action == ActionPolicyPeek {
			state.PeekedCards = []Card{state.Deck[0]}
			state.PeekerIndex = state.PresidentIndex
		}
	} else {
		state.NewTurn()
	}

	return state.EndGameIfNecessary()
}
