package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/ttodoshi/board-project/internal/core/ports"
	"github.com/ttodoshi/board-project/pkg/logging"
	"time"
)

type RoomService struct {
	redisClient *redis.Client
	log         logging.Logger
}

func NewRoomService(redisClient *redis.Client, log logging.Logger) ports.RoomService {
	return &RoomService{
		redisClient: redisClient,
		log:         log,
	}
}

func (s *RoomService) CreateRoom(userID string) (string, error) {
	result, err := s.redisClient.Get(context.Background(), fmt.Sprintf("user:%s:room", userID)).Result()
	if err == nil {
		return result, nil
	}

	roomID := uuid.NewString()
	pipe := s.redisClient.TxPipeline()
	pipe.Set(context.Background(), fmt.Sprintf("user:%s:room", userID), roomID, 24*time.Hour)
	emptyContent, _ := json.Marshal([]string{})
	pipe.HSet(context.Background(), fmt.Sprintf("room:%s", roomID),
		"admin_id", userID,
		"content", emptyContent,
	)
	pipe.Expire(context.Background(), fmt.Sprintf("room:%s", roomID), 24*time.Hour)
	_, err = pipe.Exec(context.Background())

	if err != nil {
		return "", fmt.Errorf("error while creating room: %w", err)
	}
	return roomID, nil
}
