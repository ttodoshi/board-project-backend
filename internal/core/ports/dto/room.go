package dto

type GetRoomDto struct {
	ID      string   `json:"ID,omitempty"`
	AdminID string   `json:"adminID,omitempty"`
	Content []string `json:"content,omitempty"`
}
