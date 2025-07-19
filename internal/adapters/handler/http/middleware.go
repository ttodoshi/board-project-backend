package http

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/ttodoshi/board-project/internal/core/ports/dto"
	"net/http"
	"time"
)

func ErrorHandlerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err == nil {
			return nil
		}
		var HTTPError *echo.HTTPError
		if errors.As(err, &HTTPError) {
			return c.JSON(
				HTTPError.Code,
				dto.ErrorResponseDto{
					Timestamp: time.Now(),
					Status:    HTTPError.Code,
					Error:     http.StatusText(HTTPError.Code),
					Message:   err.Error(),
					Path:      c.Request().URL.Path,
				},
			)
		}
		return c.JSON(
			http.StatusInternalServerError,
			dto.ErrorResponseDto{
				Timestamp: time.Now(),
				Status:    http.StatusInternalServerError,
				Error:     http.StatusText(http.StatusInternalServerError),
				Message:   err.Error(),
				Path:      c.Request().URL.Path,
			},
		)
	}
}
