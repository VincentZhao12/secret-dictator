package game

import (
	"math/rand"
	"time"

	"github.com/VincentZhao12/secret-hitler/backend/internal/messages"
	"github.com/VincentZhao12/secret-hitler/backend/internal/models"
)

type Game struct {
	state		models.GameState
	ID     		string
	HostID 		string
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

func NewGame() *Game {
	return &Game{
		ID: 		generateRandomID(8),
		state: 		models.NewGameState(),
	}
}

func (g *Game) SetHostID(id string) {
	g.HostID = id
}

func (g *Game) NewPlayer(username string) (*models.Player, error) {	
	player, err := g.state.AddPlayer(generateRandomID(16), username)
	if err != nil {
		return nil, err
	}

	return player, nil
}

func (g *Game) broadcastGameState() {
	// TODO: broadcast
}

func (g *Game) validateActionMessage(message messages.ActionMessage, requiresPresident bool, requiresTarget bool) *messages.ActionErrorMessage {
	if requiresPresident && g.state.Players[g.state.PresidentIndex].ID != message.SenderID {
		return messages.NewActionErrorMessage(message.SenderID, messages.NotAllowed)
	}

	if requiresTarget && message.TargetIndex < 0 || message.TargetIndex >= len(g.state.Players) {
		return messages.NewActionErrorMessage(message.SenderID, messages.InvalidTarget)
	}

	return nil

}

func (g *Game) NewTurn() {
	g.state.Phase = models.Nomination
	g.state.PrevPresidentIndex = g.state.PresidentIndex
	g.state.PrevChancellorIndex = g.state.ChancellorIndex
	
	if g.state.ResumeOrderIndex != -1 {
		g.state.PresidentIndex = g.state.ResumeOrderIndex
		g.state.ResumeOrderIndex = -1
	} else {
		g.state.PresidentIndex = (g.state.PresidentIndex + 1) % len(g.state.Players)
	}

	g.state.ChancellorIndex = -1
	g.state.NomineeIndex = -1
	g.state.Votes = nil
	g.state.PendingAction = nil
	g.state.PeekedCards = nil
	g.state.PeekerIndex = -1

	g.broadcastGameState()
}

func (g *Game) PlaceCard(card models.Card) bool {
	if card == models.CardFascist {
		g.state.Board.FascistPolicies++
	} else if card == models.CardLiberal {
		g.state.Board.LiberalPolicies++
	}

	if g.state.Board.ExecutiveActions[g.state.Board.FascistPolicies] != models.ActionNone {
		g.state.Phase = models.Executive
		action := g.state.Board.ExecutiveActions[g.state.Board.FascistPolicies]
		g.state.PendingAction = &action
		g.state.Phase = models.Executive
	}
	
	if g.state.Board.FascistPolicies >= g.state.Board.FascistSlots {
		g.state.Phase = models.GameOver
		return true
	} else if g.state.Board.LiberalPolicies >= g.state.Board.LiberalSlots {
		g.state.Phase = models.GameOver
		return true
	}
	return false
}

func (g *Game) ProcessActionMessage(message messages.ActionMessage) messages.Message {
	switch message.Action {
		case models.ActionInvestigate:
			if errorMessage := g.validateActionMessage(message, true, true); errorMessage != nil {
				return errorMessage
			}

			if g.state.Phase != models.Executive {
				return messages.NewActionErrorMessage(message.SenderID, messages.NotAllowed)
			}

			forPlayer := g.state.Players[g.state.PresidentIndex]
			player := g.state.Players[message.TargetIndex]
			
			return messages.NewGameStateMessage(
				"server",
				g.state.RevealPlayer(forPlayer, player),
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
				g.state,
			)

		case models.ActionPolicyPeek:
			if errorMessage := g.validateActionMessage(message, true, false); errorMessage != nil {
				return errorMessage
			}

			if g.state.Phase != models.Executive {
				return messages.NewActionErrorMessage(message.SenderID, messages.NotAllowed) 
			}

		case models.ActionExecution:
			if errorMessage := g.validateActionMessage(message, true, true); errorMessage != nil {
				return errorMessage
			}

			if g.state.Phase != models.Executive {
				return messages.NewActionErrorMessage(message.SenderID, messages.NotAllowed) 
			}
			
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
				g.state,
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
				g.state,
			)
		
		case models.ActionLegislate:
			if errorMessage := g.validateActionMessage(message, false, false); errorMessage != nil {
				return errorMessage
			}

			if message.SenderID != g.state.Players[g.state.PeekerIndex].ID {
				return messages.NewActionErrorMessage(message.SenderID, messages.NotAllowed) 
			}

			if message.TargetIndex < 0 || message.TargetIndex > len(g.state.PeekedCards) {
				return messages.NewActionErrorMessage(message.SenderID, messages.InvalidAction) 
			}

			if g.state.Phase == models.Legislation1 {
				g.state.Discard = append(g.state.Discard, g.state.PeekedCards[message.TargetIndex])
				g.state.PeekedCards = append(g.state.PeekedCards[:message.TargetIndex], g.state.PeekedCards[message.TargetIndex+1:]...)
				g.state.Phase = models.Legislation2
				g.state.PeekerIndex = g.state.ChancellorIndex
			} else if g.state.Phase == models.Legislation2 {
				removedCard := g.state.PeekedCards[message.TargetIndex]
				g.state.Discard = append(g.state.Discard, g.state.PeekedCards[message.TargetIndex])
				g.state.PeekedCards = nil
				g.state.PeekerIndex = -1

				if gameOver := g.PlaceCard(removedCard); !gameOver {
					g.NewTurn()
				}
			}

			g.broadcastGameState()

			return messages.NewGameStateMessage(
				"server",
				g.state,
			)

		case models.ActionProposeVeto:
			
		case models.ActionApproveVeto:
			
		case models.ActionRejectVeto:

		case models.ActionEndTurn:
			
	}

	return messages.NewActionErrorMessage(message.SenderID, messages.InvalidAction)
}