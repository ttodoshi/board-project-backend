package main

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	httpHandler "github.com/ttodoshi/board-project/internal/adapters/handler/http"
	"github.com/ttodoshi/board-project/internal/adapters/handler/http/api"
	"github.com/ttodoshi/board-project/internal/adapters/handler/ws"
	"github.com/ttodoshi/board-project/internal/core/services"
	"github.com/ttodoshi/board-project/pkg/discovery"
	"github.com/ttodoshi/board-project/pkg/env"
	"github.com/ttodoshi/board-project/pkg/logging"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	env.LoadEnvVariables()
	discovery.InitServiceDiscovery()
}

// @title		Board Project API
// @version		1.0
// @host		localhost:8080
// @BasePath	/api/v1
func main() {
	// initialize logger
	log := logging.GetLogger()

	// initialize echo
	e := echo.New()

	// initialize redis
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	// initialize services
	roomService := services.NewRoomService(rdb, log)
	connectionService := services.NewConnectionService(rdb, log)

	// redis pub/sub goroutine
	go connectionService.NotifySubscribers()

	// initialize routes
	router := httpHandler.NewRouter(
		log,
		api.NewRoomHandler(
			roomService, log,
		),
		ws.NewConnectionHandler(
			connectionService, log,
		),
	)
	router.InitRoutes(e)

	// graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := e.Start(":" + os.Getenv("PORT")); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Listen: %s\n", err)
		}
	}()
	log.Print("Server started")
	<-done
	log.Print("Server is shutting down...")

	// close redis connection
	if err := rdb.Close(); err != nil {
		log.Print("Error closing Redis connection:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Print("Server exited properly")
}
