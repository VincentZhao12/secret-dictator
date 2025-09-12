package models

import "github.com/VincentZhao12/secret-hitler/backend/internal/repository"

type ElectionTracker struct {
	FailedElections int `json:"failed_elections"`
	MaxFailures     int `json:"max_failures"`
}

type Board struct {
	LiberalPolicies 	int `json:"liberal_policies"`
	FascistPolicies 	int `json:"fascist_policies"`
	FascistSlots   		int `json:"fascist_slots"`
	LiberalSlots   		int `json:"liberal_slots"`
	ExecutiveActions    map[int]Action `json:"executive_actions"`
	ElectionTracker  	ElectionTracker `json:"election_tracker"`
}

func NewBoard(players int) (Board, error) {
	board := Board{
		LiberalPolicies: 0,
		FascistPolicies: 0,
		ExecutiveActions: make(map[int]Action),
		ElectionTracker: ElectionTracker{
			FailedElections: 0,
			MaxFailures:     3,
		},
	}

	switch players {
	case 5, 6:
		board.FascistSlots = 6
		board.LiberalSlots = 5
		board.ExecutiveActions = map[int]Action{
			3: ActionPolicyPeek,
			4: ActionExecution,
			5: ActionExecution,
		}
	case 7, 8:
		board.FascistSlots = 7
		board.LiberalSlots = 5
		board.ExecutiveActions = map[int]Action{
			2: ActionPolicyPeek,
			3: ActionExecution,
			4: ActionExecution,
			5: ActionExecution,
		}
	case 9, 10:
		board.FascistSlots = 8
		board.LiberalSlots = 5
		board.ExecutiveActions = map[int]Action{
			1: ActionInvestigate,
			2: ActionPolicyPeek,
			3: ActionExecution,
			4: ActionExecution,
			5: ActionExecution,
		}
	default:
		// TODO: Add observability for this case
		return Board{}, repository.ErrInvalidPlayerCount
	}

	return board, nil
}