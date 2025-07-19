package ports

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

var (
	BadRequestError = &echo.HTTPError{
		Internal: nil,
		Message:  "Bad Request",
		Code:     http.StatusBadRequest,
	}
	UnauthorizedError = &echo.HTTPError{
		Internal: nil,
		Message:  "Unauthorized",
		Code:     http.StatusUnauthorized,
	}
	ForbiddenError = &echo.HTTPError{
		Internal: nil,
		Message:  "Forbidden",
		Code:     http.StatusForbidden,
	}
	NotFoundError = &echo.HTTPError{
		Internal: nil,
		Message:  "Not Found",
		Code:     http.StatusNotFound,
	}
)
