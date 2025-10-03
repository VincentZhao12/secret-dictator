package game

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/VincentZhao12/secret-hitler/backend/internal/messages"
	"github.com/VincentZhao12/secret-hitler/backend/internal/models"
	"github.com/VincentZhao12/secret-hitler/backend/internal/repository"
	"github.com/gorilla/websocket"
)

type Game struct {
	state       models.GameState
	ID          string
	HostID      string
	ActionChan  chan (messages.ActionMessage)
	Connections map[string]*websocket.Conn
	manager     *Manager
	connMu      sync.RWMutex
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateRandomID(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func NewGame(manager *Manager) *Game {
	g := &Game{
		ID:          generateRandomID(8),
		state:       models.NewGameState(),
		manager:     manager,
		Connections: make(map[string]*websocket.Conn),
		ActionChan:  make(chan messages.ActionMessage),
	}
	go g.Run()
	return g
}

func (g *Game) SetHostID(id string) {
	g.HostID = id
	g.state.HostID = id
}

func (g *Game) AddConnection(id string, conn *websocket.Conn) error {
	g.connMu.Lock()
	playerIndex, exists := g.state.PlayerIndexMap[id]
	if !exists {
		fmt.Println("player doesnt exist")
		g.connMu.Unlock()
		return repository.ErrPlayerNotFound
	}

	if oldConn, exists := g.Connections[id]; exists && oldConn != nil {
		fmt.Println("player connection alr exist")
		// return repository.ErrPlayerAlreadyExists
	}

	player := g.state.GetPlayer(playerIndex)
	if player != nil {
		player.IsConnected = true
	}
	g.Connections[id] = conn

	if len(g.Connections) == len(g.state.Players) && g.state.Phase == models.Paused {
		g.state.Phase = g.state.ResumePhase
		g.state.ResumePhase = ""
	}
	playerForState := g.state.GetPlayerByID(id)
	g.connMu.Unlock()

	if playerForState != nil {
		conn.WriteJSON(messages.NewGameStateMessage(
			"server",
			g.state.ObfuscateGameState(*playerForState),
		))
	}
	g.broadcastGameState()

	return nil
}

func (g *Game) CanBeDeleted() bool {
	return len(g.Connections) == 0 && (g.state.Phase == models.GameOver || g.state.Phase == models.Setup)
}

func (g *Game) DropConnection(id string) error {
	g.connMu.Lock()
	playerIndex, exists := g.state.PlayerIndexMap[id]
	if !exists {
		g.connMu.Unlock()
		return repository.ErrPlayerNotFound
	}

	player := g.state.GetPlayer(playerIndex)
	if player != nil {
		player.IsConnected = false
	}

	delete(g.Connections, id)
	if g.state.Phase != models.GameOver && g.state.Phase != models.Setup {
		g.state.ResumePhase = g.state.Phase
		g.state.Phase = models.Paused
	}

	if g.state.Phase == models.Setup {
		err := g.state.RemovePlayer(player.ID)
		if err != nil {
			g.connMu.Unlock()
			return err
		}
		if len(g.state.Players) == 0 {
			g.connMu.Unlock()
			return nil
		}
		fmt.Println(g.HostID, id)
		if g.HostID == id {
			g.SetHostID(g.state.Players[0].ID)
		}
	}
	g.connMu.Unlock()
	g.broadcastGameState()

	return nil
}

func (g *Game) EndGame(winner models.Team) {
	g.state.EndGame(winner)
	close(g.ActionChan)
}

func (g *Game) NewPlayer(username string) (*models.Player, error) {
	playerID := generateRandomID(16)
	player, err := g.state.AddPlayer(playerID, username)
	if err != nil {
		return nil, err
	}

	// Set host to first player if not already set
	if g.HostID == "" {
		g.SetHostID(playerID)
	}

	g.broadcastGameState()

	return player, nil
}

func (g *Game) broadcastGameState() {
	g.connMu.RLock()
	snapshot := make(map[string]*websocket.Conn, len(g.Connections))
	for id, conn := range g.Connections {
		snapshot[id] = conn
	}
	g.connMu.RUnlock()

	for id, conn := range snapshot {
		if conn != nil {
			player := g.state.GetPlayerByID(id)
			if player != nil {
				if err := conn.WriteJSON(messages.NewGameStateMessage(
					"server",
					g.state.ObfuscateGameState(*player),
				)); err != nil {
					println("error sending message")
				}
			}
		}
	}
}

func (g *Game) validateActionMessage(message messages.ActionMessage, requiresPresident bool, requiresTarget bool) *messages.ActionErrorMessage {
	president := g.state.GetPlayer(g.state.PresidentIndex)
	if requiresPresident && (president == nil || president.ID != message.SenderID) {
		return messages.NewActionErrorMessage(message.SenderID, messages.NotAllowed)
	}

	if requiresTarget && (message.TargetIndex < 0 || message.TargetIndex >= len(g.state.Players)) {
		return messages.NewActionErrorMessage(message.SenderID, messages.InvalidTarget)
	}

	targetPlayer := g.state.GetPlayer(message.TargetIndex)
	if requiresTarget && (targetPlayer == nil || targetPlayer.IsExecuted) {
		return messages.NewActionErrorMessage(message.SenderID, messages.InvalidTarget)
	}

	sender := g.state.GetPlayerByID(message.SenderID)
	if sender != nil && sender.IsExecuted {
		return messages.NewActionErrorMessage(message.SenderID, messages.NotAllowed)
	}

	return nil
}

func (g *Game) NewTurn() {
	g.state.NewTurn()
	g.broadcastGameState()
}

func (g *Game) PlaceCard(card models.Card) bool {
	return g.state.PlaceCard(card)
}

func (g *Game) Run() {
	for {
		message, ok := <-g.ActionChan
		if !ok {
			return
		}
		response := g.ProcessActionMessage(message)
		g.connMu.RLock()
		conn, exists := g.Connections[message.SenderID]
		g.connMu.RUnlock()

		if !exists || conn == nil {
			continue
		}

		if err := conn.WriteJSON(response); err != nil {
			println("error writing")
			continue
		}
	}
}

// SHOULD BE IDIOMATIC
func (g *Game) ProcessActionMessage(message messages.ActionMessage) messages.Message {
	p := g.state.GetPlayerByID(message.SenderID)
	if p == nil {
		return messages.NewActionErrorMessage(message.SenderID, messages.NotAllowed)
	}
	switch message.Action {
	case models.ActionStartGame:
		// Only host can start the game
		if message.SenderID != g.HostID {
			return messages.NewActionErrorMessage(message.SenderID, messages.NotAllowed)
		}

		err := g.state.StartGame()
		if err != nil {
			return messages.NewActionErrorMessage(message.SenderID, messages.CouldNotStart)
		}

		g.broadcastGameState()

		return messages.NewGameStateMessage(
			"server",
			g.state.ObfuscateGameState(*p),
		)

	case models.ActionInvestigate:
		if errorMessage := g.validateActionMessage(message, true, true); errorMessage != nil {
			return errorMessage
		}

		if g.state.Phase != models.Executive {
			return messages.NewActionErrorMessage(message.SenderID, messages.NotAllowed)
		}

		forPlayer := g.state.GetPlayer(g.state.PresidentIndex)
		player := g.state.GetPlayer(message.TargetIndex)

		if forPlayer == nil || player == nil {
			return messages.NewActionErrorMessage(message.SenderID, messages.InvalidAction)
		}

		return messages.NewGameStateMessage(
			"server",
			g.state.RevealPlayer(*forPlayer, *player),
		)

	case models.ActionSpecialElection:
		if errorMessage := g.validateActionMessage(message, true, true); errorMessage != nil {
			return errorMessage
		}

		if g.state.Phase != models.Executive {
			return messages.NewActionErrorMessage(message.SenderID, messages.NotAllowed)
		}

		g.state.ResumeOrderIndex = (g.state.PresidentIndex + 1) % len(g.state.Players)

		newPresidentIndex := message.TargetIndex
		g.state.PresidentIndex = newPresidentIndex
		g.state.Phase = models.Nomination
		g.state.ChancellorIndex = -1
		g.state.NomineeIndex = -1
		g.state.Votes = nil
		g.state.PendingAction = nil
		g.state.PeekedCards = nil
		g.state.PeekerIndex = -1

		g.broadcastGameState()

		return messages.NewGameStateMessage(
			"server",
			g.state.ObfuscateGameState(*p),
		)

	case models.ActionPolicyPeek:
		if errorMessage := g.validateActionMessage(message, true, false); errorMessage != nil {
			return errorMessage
		}

		if g.state.Phase != models.Executive {
			return messages.NewActionErrorMessage(message.SenderID, messages.NotAllowed)
		}
		g.state.PeekedCards = g.state.Deck[:3]

		return messages.NewGameStateMessage(
			"server",
			g.state.ObfuscateGameState(*p),
		)

	case models.ActionExecution:
		if errorMessage := g.validateActionMessage(message, true, true); errorMessage != nil {
			return errorMessage
		}

		if g.state.Phase != models.Executive {
			return messages.NewActionErrorMessage(message.SenderID, messages.NotAllowed)
		}

		targetPlayer := g.state.GetPlayer(message.TargetIndex)
		if targetPlayer != nil {
			targetPlayer.IsExecuted = true
		}

		if g.state.EndGameIfNecessary() {
			g.broadcastGameState()
		} else {
			g.NewTurn()
		}

		return messages.NewGameStateMessage(
			"server",
			g.state.ObfuscateGameState(*p),
		)

	case models.ActionVote:
		if errorMessage := g.validateActionMessage(message, false, true); errorMessage != nil {
			return errorMessage
		}

		if g.state.Phase != models.Election {
			return messages.NewActionErrorMessage(message.SenderID, messages.NotAllowed)
		}

		if message.Vote == nil {
			return messages.NewActionErrorMessage(message.SenderID, messages.InvalidAction)
		}

		if *message.Vote {
			g.state.Votes[g.state.PlayerIndexMap[message.SenderID]] = models.VoteJa
		} else {
			g.state.Votes[g.state.PlayerIndexMap[message.SenderID]] = models.VoteNein
		}

		votes := 0
		yesVotes := 0
		for _, vote := range g.state.Votes {
			if vote != models.VotePending {
				votes++
			}
			if vote == models.VoteJa {
				yesVotes++
			}
		}

		if votes == len(g.state.Players) {
			if yesVotes > len(g.state.Players)/2 {
				g.state.ChancellorIndex = g.state.NomineeIndex
				g.state.Phase = models.Legislation1
				g.state.PeekerIndex = g.state.PresidentIndex

				fmt.Println(g.state.Deck)

				// Draw 3 cards
				if len(g.state.Deck) < 3 {
					g.state.ShuffleDeck()
				}
				g.state.PeekedCards = g.state.Deck[:3]
				g.state.Deck = g.state.Deck[3:]
			} else {
				g.state.Board.ElectionTracker.FailedElections++

				if g.state.Board.ElectionTracker.FailedElections == g.state.Board.ElectionTracker.MaxFailures {
					g.state.Board.ElectionTracker.FailedElections = 0
					// Draw top policy
					if len(g.state.Deck) < 1 {
						g.state.ShuffleDeck()
					}
					topCard := g.state.Deck[0]
					g.state.Deck = g.state.Deck[1:]

					g.PlaceCard(topCard)
				}
				g.NewTurn()
			}
		}

		g.broadcastGameState()

		return messages.NewGameStateMessage(
			"server",
			g.state.ObfuscateGameState(*p),
		)

	case models.ActionNominate:
		if errorMessage := g.validateActionMessage(message, true, true); errorMessage != nil {
			return errorMessage
		}

		if g.state.Phase != models.Nomination {
			return messages.NewActionErrorMessage(message.SenderID, messages.NotAllowed)
		}

		if message.TargetIndex == g.state.PresidentIndex ||
			message.TargetIndex == g.state.PrevPresidentIndex ||
			message.TargetIndex == g.state.PrevChancellorIndex {

			return messages.NewActionErrorMessage(message.SenderID, messages.InvalidTarget)
		}

		g.state.NomineeIndex = message.TargetIndex
		g.state.Phase = models.Election
		g.state.Votes = make([]models.VoteResult, len(g.state.Players))

		g.broadcastGameState()

		return messages.NewGameStateMessage(
			"server",
			g.state.ObfuscateGameState(*p),
		)

	case models.ActionLegislate:
		if errorMessage := g.validateActionMessage(message, false, false); errorMessage != nil {
			return errorMessage
		}

		peeker := g.state.GetPlayer(g.state.PeekerIndex)
		if peeker == nil || message.SenderID != peeker.ID {
			return messages.NewActionErrorMessage(message.SenderID, messages.NotAllowed)
		}

		if message.TargetIndex < 0 || message.TargetIndex > len(g.state.PeekedCards) {
			return messages.NewActionErrorMessage(message.SenderID, messages.InvalidAction)
		}

		switch g.state.Phase {
		case models.Legislation1:
			g.state.Discard = append(g.state.Discard, g.state.PeekedCards[message.TargetIndex])
			g.state.PeekedCards = append(g.state.PeekedCards[:message.TargetIndex], g.state.PeekedCards[message.TargetIndex+1:]...)
			g.state.Phase = models.Legislation2
			g.state.PeekerIndex = g.state.ChancellorIndex
		case models.Legislation2:
			removedCard := g.state.PeekedCards[message.TargetIndex]
			g.state.Discard = append(g.state.Discard, g.state.PeekedCards[message.TargetIndex])
			g.state.PeekedCards = nil
			g.state.PeekerIndex = -1

			if !g.PlaceCard(removedCard) {
				g.NewTurn()
			}
		}

		g.broadcastGameState()

		return messages.NewGameStateMessage(
			"server",
			g.state.ObfuscateGameState(*p),
		)

	case models.ActionProposeVeto:

	case models.ActionApproveVeto:

	case models.ActionRejectVeto:

	case models.ActionEndTurn:
		if errorMessage := g.validateActionMessage(message, false, false); errorMessage != nil {
			return errorMessage
		}

		if g.state.Phase != models.Executive {
			return messages.NewActionErrorMessage(message.SenderID, messages.NotAllowed)
		}

		g.NewTurn()

		return messages.NewGameStateMessage(
			"server",
			g.state.ObfuscateGameState(*p),
		)
	}

	return messages.NewActionErrorMessage(message.SenderID, messages.InvalidAction)
}
