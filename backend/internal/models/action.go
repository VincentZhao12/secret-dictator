package models

type Action string

const (
	ActionStartGame       Action = "start_game"
	ActionInvestigate     Action = "investigate"
	ActionSpecialElection Action = "special_election"
	ActionPolicyPeek      Action = "policy_peek"
	ActionExecution       Action = "execution"
	ActionVote            Action = "vote"
	ActionNominate        Action = "nominate"
	ActionLegislate       Action = "legislate"
	ActionProposeVeto     Action = "propose_veto"
	ActionApproveVeto     Action = "approve_veto"
	ActionRejectVeto      Action = "reject_veto"
	ActionEndTurn         Action = "end_turn"
	ActionNone            Action = "none"
)
