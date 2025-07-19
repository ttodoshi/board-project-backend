package dto

import "time"

type ErrorResponseDto struct {
	Timestamp time.Time `json:"timestamp"`
	Status    int       `json:"status,omitempty"`
	Error     string    `json:"error,omitempty"`
	Message   string    `json:"message,omitempty"`
	Path      string    `json:"path,omitempty"`
}
