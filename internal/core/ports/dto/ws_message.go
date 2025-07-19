package dto

type MessageType string

const (
	Update MessageType = "update"
	Join   MessageType = "join"
)

type ClientMessage struct {
	MessageType MessageType `json:"type"`
	RoomID      string      `json:"roomID"`
	Content     []string    `json:"content"`
}
