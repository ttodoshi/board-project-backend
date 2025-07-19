package api

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ttodoshi/board-project/internal/core/ports"
	"github.com/ttodoshi/board-project/pkg/logging"
	"net/http"
)

type RoomHandler struct {
	svc ports.RoomService
	log logging.Logger
}

func NewRoomHandler(svc ports.RoomService, log logging.Logger) *RoomHandler {
	return &RoomHandler{
		svc: svc,
		log: log,
	}
}

// CreateRoom godoc
//
//	@Summary		Create room
//	@Description	Create room
//	@Tags			rooms
//	@Accept			json
//	@Produce		json
//	@Param			Cookie	header		string	true	"userID"	default(userID=)
//	@Success		201		{object}	string
//	@Failure		404		{object}	dto.ErrorResponseDto
//	@Router			/rooms [post]
func (h *RoomHandler) CreateRoom(c echo.Context) error {
	h.log.Debug("received create room request")

	var userID string
	if userIDCookie, err := c.Cookie("userID"); err == nil && userIDCookie.Value != "" {
		userID = userIDCookie.Value
	} else {
		userID = uuid.NewString()
		c.SetCookie(&http.Cookie{
			Name:     "userID",
			Value:    userID,
			Path:     "/",
			HttpOnly: true,
		})
	}
	ID, err := h.svc.CreateRoom(userID)
	if err != nil {
		return err
	}

	return c.String(http.StatusCreated, ID)
}
