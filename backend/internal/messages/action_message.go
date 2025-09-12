package messages

import "github.com/VincentZhao12/secret-hitler/backend/internal/models"

const (
	MesasgeTypeAction MessageType = "action"
	MesasgeTypeActionError MessageType = "action_error"
)


type ActionMessage struct {
	BaseMessage
	Action			models.Action `json:"action"`
	TargetIndex		int `json:"target_index"`
}

func NewActionMessage(senderID string, action models.Action, targetIndex int  ) *ActionMessage {
	return &ActionMessage{
		BaseMessage: BaseMessage{
			Type:     MesasgeTypeAction,
			SenderID: senderID,
		},
		Action: action,
		TargetIndex: targetIndex,
	}
}

type ActionErrorReason string

const (
	NotAllowed ActionErrorReason = "Action not allowed"
	InvalidTarget ActionErrorReason = "Targeted player is cannot be targeted for this action"
	InvalidAction ActionErrorReason = "Invalid action"
)

type ActionErrorMessage struct {
	BaseMessage
	Reason		ActionErrorReason `json:"target_index"`
}

func NewActionErrorMessage(senderID string, reason ActionErrorReason) *ActionErrorMessage {
	return &ActionErrorMessage{
		BaseMessage: BaseMessage{
			Type:     MesasgeTypeActionError,
			SenderID: senderID,
		},
		Reason: reason,
	}
}