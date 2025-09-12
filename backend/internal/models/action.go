package models

type Action string

const (
	ActionInvestigate 		Action = "investigate"
	ActionSpecialElection 	Action = "special_election"
	ActionPolicyPeek 		Action = "policy_peek"
	ActionExecution 		Action = "execution"
	ActionPass 				Action = "pass"
	ActionNominate 			Action = "nominate"
	ActionProposeVeto 		Action = "propose_veto"
	ActionApproveVeto 		Action = "approve_veto"
	ActionRejectVeto 		Action = "reject_veto"
)