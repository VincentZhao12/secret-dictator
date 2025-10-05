package models

import "github.com/VincentZhao12/secret-hitler/backend/internal/repository"

type ElectionTracker struct {
	FailedElections int `json:"failed_elections"`
	MaxFailures     int `json:"max_failures"`
}

type Board struct {
	LiberalPolicies  int             `json:"liberal_policies"`
	FascistPolicies  int             `json:"fascist_policies"`
	FascistSlots     int             `json:"fascist_slots"`
	LiberalSlots     int             `json:"liberal_slots"`
	ExecutiveActions map[int]Action  `json:"executive_actions"`
	ElectionTracker  ElectionTracker `json:"election_tracker"`
	DangerZoneStart  int             `json:"danger_zone_start"`
}

func NewBoard(players int) (Board, error) {
	b := Board{
		LiberalPolicies:  0,
		FascistPolicies:  0,
		ExecutiveActions: make(map[int]Action),
		ElectionTracker: ElectionTracker{
			FailedElections: 0,
			MaxFailures:     3,
		},
		// liberal slots always 5
		LiberalSlots: 5,
	}

	switch players {
	case 5, 6:
		b.FascistSlots = 6
		// From the board: the fascist track with 6 slots has powers at slot 3 (peek) and slot 4 (execution) etc.
		b.ExecutiveActions = map[int]Action{
			1: ActionSpecialElection,
			3: ActionPolicyPeek,
			4: ActionExecution,
			5: ActionExecution,
		}
		// For 5–6 players, the danger zone (i.e. when fascist is two away from winning) starts earlier
		b.DangerZoneStart = 3
	case 7, 8:
		b.FascistSlots = 7
		b.ExecutiveActions = map[int]Action{
			2: ActionPolicyPeek,
			3: ActionExecution,
			4: ActionExecution,
			5: ActionExecution,
		}
		b.DangerZoneStart = 4
	case 9, 10:
		b.FascistSlots = 8
		b.ExecutiveActions = map[int]Action{
			1: ActionInvestigate,
			2: ActionSpecialElection,
			3: ActionExecution,
			4: ActionExecution,
			5: ActionExecution,
		}
		b.DangerZoneStart = 5
	default:
		return Board{}, repository.ErrInvalidPlayerCount
	}

	return b, nil
}
