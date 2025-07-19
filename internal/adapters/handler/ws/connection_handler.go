package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/ttodoshi/board-project/internal/core/ports"
	"github.com/ttodoshi/board-project/internal/core/ports/dto"
	"github.com/ttodoshi/board-project/pkg/logging"
	"log"
	"net/http"
)

type ConnectionHandler struct {
	svc      ports.ConnectionService
	log      logging.Logger
	upgrader websocket.Upgrader
}

func NewConnectionHandler(svc ports.ConnectionService, log logging.Logger) *ConnectionHandler {
	return &ConnectionHandler{
		svc: svc,
		log: log,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

func (h *ConnectionHandler) HandleWebSocket(c echo.Context) error {
	var userID string
	if cookie, err := c.Cookie("userID"); err != nil {
		userID = ""
	} else {
		userID = cookie.Value
	}

	conn, err := h.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println("Upgrade:", err)
		return err
	}
	defer func(conn *websocket.Conn) {
		err = conn.Close()
		if err != nil {
			h.log.Error(err)
		}
	}(conn)
	defer func(svc ports.ConnectionService) {
		err = svc.CloseConnections()
		if err != nil {
			h.log.Error(err)
		}
	}(h.svc)

	for {
		var msg dto.ClientMessage
		err = conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				return err
			}
			log.Println("Read error:", err)
			continue
		}
		switch msg.MessageType {
		case dto.Update:
			err = h.svc.UpdateRoom(msg, userID)
		case dto.Join:
			err = h.svc.JoinRoom(msg.RoomID, conn)
		default:
			err = fmt.Errorf("unsupported message type: %s", msg.MessageType)
		}
		if err != nil {
			err = conn.WriteJSON(map[string]interface{}{
				"message": err.Error(),
			})
		}
		if err != nil {
			log.Println("Write:", err)
		}
	}
}
