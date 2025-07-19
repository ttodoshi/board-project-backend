package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"github.com/ttodoshi/board-project/internal/core/ports"
	"github.com/ttodoshi/board-project/internal/core/ports/dto"
	"github.com/ttodoshi/board-project/pkg/logging"
	"log"
	"sync"
)

type ConnectionService struct {
	mu          sync.Mutex
	rooms       map[string]map[*websocket.Conn]struct{}
	redisClient *redis.Client
	pubsub      *redis.PubSub
	log         logging.Logger
}

func NewConnectionService(redisClient *redis.Client, log logging.Logger) ports.ConnectionService {
	return &ConnectionService{
		rooms:       map[string]map[*websocket.Conn]struct{}{},
		redisClient: redisClient,
		pubsub:      redisClient.Subscribe(context.Background()),
		log:         log,
	}
}

func (s *ConnectionService) NotifySubscribers() {
	ch := s.pubsub.Channel()

	for message := range ch {
		var msg []string
		if err := json.Unmarshal([]byte(message.Payload), &msg); err != nil {
			log.Println("JSON unmarshal error:", err)
			continue
		}
		s.mu.Lock()
		for client := range s.rooms[message.Channel] {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Write error: %v", err)
				err = client.Close()
				if err != nil {
					log.Printf("Close error: %v", err)
					continue
				}
				delete(s.rooms[message.Channel], client)
			}
		}
		s.mu.Unlock()
	}
}

func (s *ConnectionService) UpdateRoom(msg dto.ClientMessage, userID string) error {
	adminID, err := s.redisClient.HGet(context.Background(), fmt.Sprintf("room:%s", msg.RoomID), "admin_id").Result()
	if err != nil {
		return err
	}
	if adminID != userID {
		return fmt.Errorf("you're not admin of this room")
	}
	var content []byte
	content, err = json.Marshal(msg.Content)
	if err != nil {
		return err
	}
	s.redisClient.HSet(context.Background(), fmt.Sprintf("room:%s", msg.RoomID), "content", content)
	s.redisClient.Publish(context.Background(), fmt.Sprintf("room:%s", msg.RoomID), content)
	return nil
}

func (s *ConnectionService) JoinRoom(roomID string, conn *websocket.Conn) error {
	roomKey := fmt.Sprintf("room:%s", roomID)
	s.mu.Lock()
	if _, ok := s.rooms[roomKey]; !ok {
		s.rooms[roomKey] = map[*websocket.Conn]struct{}{}
	}
	s.rooms[roomKey][conn] = struct{}{}
	s.mu.Unlock()

	content, err := s.redisClient.HGet(context.Background(), fmt.Sprintf("room:%s", roomID), "content").Result()
	if err != nil {
		return err
	}
	err = conn.WriteJSON(content)
	if err != nil {
		return err
	}
	return s.pubsub.Subscribe(context.Background(), fmt.Sprintf("room:%s", roomID))
}

func (s *ConnectionService) CloseConnections() error {
	s.mu.Lock()
	for _, room := range s.rooms {
		for conn := range room {
			err := conn.Close()
			if err != nil {
				return err
			}
		}
	}
	s.rooms = map[string]map[*websocket.Conn]struct{}{}
	s.mu.Unlock()
	return s.pubsub.PUnsubscribe(context.Background(), "room:*")
}
