package messages

const (
	MessageTypeConnectionError MessageType = "connection_error"
)

type ConnectionErrorType int

const (
	ConnectionErrorTypeFatal ConnectionErrorType = iota
	ConnectionErrorTypeGameInvalid
	ConnectionErrorTypePlayerInvalid
	ConnectionErrorTypeServerError
)

type ConnectionErrorMessage struct {
	BaseMessage `json:"base_message" tstype:"BaseMessage"`
	Reason      string              `json:"reason"`
	ErrorType   ConnectionErrorType `json:"error_type"`
}

func NewConnectionErrorMessage(senderID string, reason string, errorType ConnectionErrorType) *ConnectionErrorMessage {
	return &ConnectionErrorMessage{
		BaseMessage: BaseMessage{
			Type:     MessageTypeConnectionError,
			SenderID: senderID,
		},
		Reason:    reason,
		ErrorType: errorType,
	}
}
