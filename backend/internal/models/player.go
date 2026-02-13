package models

type PlayerRole string

const (
	RoleLiberal    PlayerRole = "liberal"
	RoleHitler     PlayerRole = "hitler"
	RoleFascist    PlayerRole = "fascist"
	RoleHidden     PlayerRole = "hidden"
	RoleUnassigned PlayerRole = "unassigned"
)

type Player struct {
	ID          string     `json:"id"`
	Username    string     `json:"username"`
	Role        PlayerRole `json:"role"`
	IsExecuted  bool       `json:"is_executed"`
	IsConnected bool       `json:"is_connected"`
	IsBot       bool       `json:"is_bot"`
	ModelSlug   string     `json:"model_slug"`
}

var BotUsernames = []string{
	"Robespierre",
	"Nelson Mandela",
	"Optimus Prime",
	"Jor Jor Well",
	"Tyrion Lannister",
	"Luthen Rael",
	"Jar Jar Binks",
	"Winston Churchill",
	"Petyr Baelish",
	"Fidel Castro",
	"Che Guevara",
	"Sun Tzu",
	"George Washington",
	"Cassian Andor",
	"Mon Mothma",
	"Saw Gerrera",
	"Sheev Palpatine",
	"Orson Krennic",
	"Benedict Arnold",
	"Ho Chi Minh",
	"Stilgar",
	"Paul Atreides",
	"Leto Atreides",
}

func NewPlayer(id string, username string) Player {
	return Player{
		ID:         id,
		Username:   username,
		Role:       RoleUnassigned,
		IsExecuted: false,
		IsBot:      false,
		ModelSlug:  "",
	}
}

func NewBotPlayer(id string, modelSlug string, username string) Player {
	return Player{
		ID:         id,
		Username:   username,
		Role:       RoleUnassigned,
		IsExecuted: false,
		IsBot:      true,
		ModelSlug:  modelSlug,
	}
}