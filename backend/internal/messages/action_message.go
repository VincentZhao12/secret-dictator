package messages

import "github.com/VincentZhao12/secret-hitler/backend/internal/models"

const (
	MessageTypeAction      MessageType = "action"
	MessageTypeActionError MessageType = "action_error"
)

type ActionMessage struct {
	BaseMessage `json:"base_message" tstype:"BaseMessage"`
	Action      models.Action `json:"action" tstype:"Action"`
	TargetIndex int           `json:"target_index,omitempty"`
	Vote        *bool         `json:"vote,omitempty"`
}

func NewActionMessage(senderID string, action models.Action, targetIndex int, vote bool) *ActionMessage {
	return &ActionMessage{
		BaseMessage: BaseMessage{
			Type:     MessageTypeAction,
			SenderID: senderID,
		},
		Action:      action,
		TargetIndex: targetIndex,
		Vote:        &vote,
	}
}

type ActionErrorReason string

const (
	NotAllowed    ActionErrorReason = "Action not allowed"
	InvalidTarget ActionErrorReason = "Targeted player is cannot be targeted for this action"
	InvalidAction ActionErrorReason = "Invalid action"
	CouldNotStart ActionErrorReason = "Could not start game"
)

type ActionErrorMessage struct {
	BaseMessage `json:"base_message" tstype:"BaseMessage"`
	Reason ActionErrorReason `json:"reason"`
}

func NewActionErrorMessage(senderID string, reason ActionErrorReason) *ActionErrorMessage {
	return &ActionErrorMessage{
		BaseMessage: BaseMessage{
			Type:     MessageTypeActionError,
			SenderID: senderID,
		},
		Reason: reason,
	}
}
