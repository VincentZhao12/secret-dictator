package models

type Action string

const (
	ActionInvestigate 		Action = "investigate"
	ActionSpecialElection 	Action = "special_election"
	ActionPolicyPeek 		Action = "policy_peek"
	ActionExecution 		Action = "execution"
	ActionPass 				Action = "pass"
)