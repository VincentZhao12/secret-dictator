package models

type ElectionTracker struct {
	FailedElections int `json:"failed_elections"`
	MaxFailures     int `json:"max_failures"`
}

type Board struct {
	LiberalPolicies 	int `json:"liberal_policies"`
	FascistPolicies 	int `json:"fascist_policies"`
	FascistSlots   		int `json:"fascist_slots"`
	LiberalSlots   		int `json:"liberal_slots"`
	Actions       		map[int]Action `json:"actions"`
	ElectionTracker  	ElectionTracker `json:"election_tracker"`
}