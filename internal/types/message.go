package types

import "time"

type MessageType string

const (
	MessageTypeUser      MessageType = "user"
	MessageTypeAssistant MessageType = "assistant"
	MessageTypeSystem    MessageType = "system"
)

type Message struct {
	UUID      string         `json:"uuid"`
	Type      MessageType    `json:"type"`
	Role      string         `json:"role,omitempty"`
	Content   []ContentBlock `json:"content,omitempty"`
	Usage     *Usage         `json:"usage,omitempty"`
	Timestamp time.Time      `json:"timestamp"`

	StopReason string `json:"stop_reason,omitempty"`
}

func NewUserMessage(content string) *Message {
	return &Message{
		UUID:      generateUUID(),
		Type:      MessageTypeUser,
		Role:      "user",
		Content:   []ContentBlock{{Type: "text", Text: content}},
		Timestamp: time.Now(),
	}
}

func NewAssistantMessage() *Message {
	return &Message{
		UUID:      generateUUID(),
		Type:      MessageTypeAssistant,
		Role:      "assistant",
		Timestamp: time.Now(),
	}
}
