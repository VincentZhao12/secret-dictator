package messages

type MessageType string

type Message interface {
	GetType() MessageType
	GetSenderID() string
}

type BaseMessage struct {
	Type     MessageType `json:"type"`
	SenderID string      `json:"sender_id"`
}

func (m *BaseMessage) GetType() MessageType {
	return m.Type
}

func (m *BaseMessage) GetSenderID() string {
	return m.SenderID
}
