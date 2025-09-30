package models

type PlayerRole string

const (
	RoleLiberal 	PlayerRole = "liberal"
	RoleHitler  	PlayerRole = "hitler"
	RoleFascist 	PlayerRole = "fascist"
	RoleHidden  	PlayerRole = "hidden"
	RoleUnassigned 	PlayerRole = "unassigned"
)


type Player struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     PlayerRole `json:"role"`
	IsExecuted bool `json:"is_executed"`
}

func NewPlayer(id string, username string) Player {
	return Player{
		ID:       id,
		Username: username,
		Role:     RoleUnassigned,
		IsExecuted: false,
	}
}