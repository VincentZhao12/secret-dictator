package models

type PlayerRole string

const (
	RoleLiberal PlayerRole = "liberal"
	RoleHitler  PlayerRole = "hitler"
	RoleFascist PlayerRole = "fascist"
	RoleHidden  PlayerRole = "hidden"
)


type Player struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     PlayerRole `json:"role"`
}

