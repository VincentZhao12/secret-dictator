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
}

func NewPlayer(id string, username string, role PlayerRole) Player {
	return Player{
		ID:       id,
		Username: username,
		Role:     role,
	}
}