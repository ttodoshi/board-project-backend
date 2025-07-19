package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "github.com/ttodoshi/board-project/docs"
	"github.com/ttodoshi/board-project/internal/adapters/handler/http/api"
	"github.com/ttodoshi/board-project/internal/adapters/handler/ws"
	"github.com/ttodoshi/board-project/pkg/logging"
	"net/http"
)

type Router struct {
	log logging.Logger
	*api.RoomHandler
	*ws.ConnectionHandler
}

func NewRouter(log logging.Logger, roomHandler *api.RoomHandler, connectionHandler *ws.ConnectionHandler) *Router {
	return &Router{
		log:               log,
		RoomHandler:       roomHandler,
		ConnectionHandler: connectionHandler,
	}
}

func (r *Router) InitRoutes(e *echo.Echo) {
	r.log.Info("initializing middleware")
	e.Use(ErrorHandlerMiddleware)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOriginFunc: func(origin string) (bool, error) {
			return origin != "", nil
		},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodOptions, http.MethodHead},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.log.Info("initializing routes")

	// swagger
	e.GET("/swagger-ui/*any", echoSwagger.WrapHandler)

	// healthcheck
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	apiGroup := e.Group("/api")

	v1ApiGroup := apiGroup.Group("/v1")
	{
		v1ApiGroup.GET("/ws", r.HandleWebSocket)
		v1ApiGroup.POST("/rooms", r.CreateRoom)
	}
}
