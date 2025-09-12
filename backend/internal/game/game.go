package game

import (
	"time"

	"math/rand"

	"github.com/VincentZhao12/secret-hitler/backend/internal/models"
)

type Game struct {
	state		models.GameState
	ID     		string
	HostID 		string
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateRandomGameID(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func NewGame() *Game {
	return &Game{
		ID:      generateRandomGameID(8),
	}
}

func (g *Game) SetHostID(id string) {
	g.HostID = id
}

func (g *Game) AddPlayer(id string) {
	
}




