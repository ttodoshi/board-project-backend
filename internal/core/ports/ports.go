package ports

import (
	"github.com/gorilla/websocket"
	"github.com/ttodoshi/board-project/internal/core/ports/dto"
)

type RoomService interface {
	CreateRoom(userID string) (string, error)
}

type ConnectionService interface {
	JoinRoom(roomID string, conn *websocket.Conn) error
	UpdateRoom(msg dto.ClientMessage, userID string) error
	CloseConnections() error
	NotifySubscribers()
}
